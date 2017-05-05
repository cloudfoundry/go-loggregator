package loggregator_v2

import (
	"context"
	"time"

	"code.cloudfoundry.org/lager"

	"github.com/cloudfoundry/sonde-go/events"
)

//go:generate counterfeiter -o fakes/fake_ingress_server.go . IngressServer
//go:generate counterfeiter -o fakes/fake_ingress_sender_server.go . Ingress_SenderServer

type envelopeWithResponseChannel struct {
	envelope *Envelope
	errCh    chan error
}

type Connector func() (IngressClient, error)

type grpcClient struct {
	logger           lager.Logger
	ingressClient    IngressClient
	sender           Ingress_SenderClient
	batchSender      Ingress_BatchSenderClient
	envelopes        chan *envelopeWithResponseChannel
	batchedEnvelopes chan *envelopeWithResponseChannel
	connector        Connector
	config           *MetronConfig
}

func newGrpcClient(
	logger lager.Logger,
	config *MetronConfig,
	ingressClient IngressClient,
) *grpcClient {
	client := &grpcClient{
		logger:           logger.Session("grpc-client"),
		ingressClient:    ingressClient,
		config:           config,
		envelopes:        make(chan *envelopeWithResponseChannel),
		batchedEnvelopes: make(chan *envelopeWithResponseChannel),
	}

	go client.startSender()
	go client.startBatchSender()

	return client
}

func (c *grpcClient) startSender() {
	for {
		envelopeWithResponseChannel := <-c.envelopes
		envelope := envelopeWithResponseChannel.envelope
		errCh := envelopeWithResponseChannel.errCh
		if c.sender == nil {
			var err error
			c.sender, err = c.ingressClient.Sender(context.Background())
			if err != nil {
				c.logger.Error("failed-to-create-grpc-sender", err)
				errCh <- err
				continue
			}
		}
		err := c.sender.Send(envelope)
		if err != nil {
			c.sender = nil
		}
		errCh <- err
	}
}

func (c *grpcClient) startBatchSender() {
	for {
		envelopeWithResponseChannel := <-c.batchedEnvelopes
		envelope := envelopeWithResponseChannel.envelope
		errCh := envelopeWithResponseChannel.errCh
		if c.batchSender == nil {
			var err error
			c.batchSender, err = c.ingressClient.BatchSender(context.Background())
			if err != nil {
				c.logger.Error("failed-to-create-grpc-sender", err)
				errCh <- err
				continue
			}
		}
		err := c.batchSender.Send(&EnvelopeBatch{Batch: []*Envelope{envelope}})
		if err != nil {
			c.batchSender = nil
		}
		errCh <- err
	}
}

func newTextValue(t string) *Value {
	return &Value{Data: &Value_Text{Text: t}}
}

func newGaugeValue(f float64) *GaugeValue {
	return &GaugeValue{Value: f}
}

func newGaugeValueFromUInt64(i uint64) *GaugeValue {
	return newGaugeValue(float64(i))
}

func (c *grpcClient) addEnvelopeTags(env *Envelope) {
	if env.Tags == nil {
		env.Tags = make(map[string]*Value)
	}
	env.Tags["deployment"] = newTextValue(c.config.JobDeployment)
	env.Tags["job"] = newTextValue(c.config.JobName)
	env.Tags["index"] = newTextValue(c.config.JobIndex)
	env.Tags["ip"] = newTextValue(c.config.JobIP)
	env.Tags["origin"] = newTextValue(c.config.JobOrigin)
}

func (c *grpcClient) createLogEnvelope(appID, message, sourceType, sourceInstance string, logType Log_Type) *Envelope {
	env := &Envelope{
		Timestamp:  time.Now().UnixNano(),
		SourceId:   appID,
		InstanceId: sourceInstance,
		Message: &Envelope_Log{
			Log: &Log{
				Payload: []byte(message),
				Type:    logType,
			},
		},
		Tags: map[string]*Value{
			"source_type":     newTextValue(sourceType),
			"source_instance": newTextValue(sourceInstance),
		},
	}
	return env
}

func (c *grpcClient) send(envelope *Envelope) error {
	c.addEnvelopeTags(envelope)

	e := &envelopeWithResponseChannel{
		envelope: envelope,
		errCh:    make(chan error),
	}
	defer close(e.errCh)

	c.envelopes <- e
	err := <-e.errCh
	return err
}

func (c *grpcClient) sendInBatch(envelope *Envelope) error {
	c.addEnvelopeTags(envelope)

	e := &envelopeWithResponseChannel{
		envelope: envelope,
		errCh:    make(chan error),
	}
	defer close(e.errCh)

	c.batchedEnvelopes <- e
	err := <-e.errCh
	return err
}

func (c *grpcClient) Batcher() Batcher {
	return &grpcBatcher{
		c:       c,
		metrics: make(map[string]*GaugeValue),
	}
}

func (c *grpcClient) SendAppLog(appID, message, sourceType, sourceInstance string) error {
	return c.sendInBatch(c.createLogEnvelope(appID, message, sourceType, sourceInstance, Log_OUT))
}

func (c *grpcClient) SendAppErrorLog(appID, message, sourceType, sourceInstance string) error {
	return c.send(c.createLogEnvelope(appID, message, sourceType, sourceInstance, Log_ERR))
}

func (c *grpcClient) SendAppMetrics(m *events.ContainerMetric) error {
	env := &Envelope{
		Timestamp: time.Now().UnixNano(),
		SourceId:  m.GetApplicationId(),
		Message: &Envelope_Gauge{
			Gauge: &Gauge{
				Metrics: map[string]*GaugeValue{
					"instance_index": newGaugeValue(float64(m.GetInstanceIndex())),
					"cpu":            newGaugeValue(m.GetCpuPercentage()),
					"memory":         newGaugeValueFromUInt64(m.GetMemoryBytes()),
					"disk":           newGaugeValueFromUInt64(m.GetDiskBytes()),
					"memory_quota":   newGaugeValueFromUInt64(m.GetMemoryBytesQuota()),
					"disk_quota":     newGaugeValueFromUInt64(m.GetDiskBytesQuota()),
				},
			},
		},
	}
	return c.send(env)
}

func (c *grpcClient) SendDuration(name string, duration time.Duration) error {
	b := c.Batcher()
	b.SendDuration(name, duration)
	return b.Send()
}

func (c *grpcClient) SendMebiBytes(name string, mebibytes int) error {
	b := c.Batcher()
	b.SendMebiBytes(name, mebibytes)
	return b.Send()
}

func (c *grpcClient) SendMetric(name string, value int) error {
	b := c.Batcher()
	b.SendMetric(name, value)
	return b.Send()
}

func (c *grpcClient) SendBytesPerSecond(name string, value float64) error {
	b := c.Batcher()
	b.SendBytesPerSecond(name, value)
	return b.Send()
}

func (c *grpcClient) SendRequestsPerSecond(name string, value float64) error {
	b := c.Batcher()
	b.SendRequestsPerSecond(name, value)
	return b.Send()
}

func (c *grpcClient) IncrementCounter(name string) error {
	env := &Envelope{
		Timestamp: time.Now().UnixNano(),
		Message: &Envelope_Counter{
			Counter: &Counter{
				Name: name,
				Value: &Counter_Delta{
					Delta: uint64(1),
				},
			},
		},
	}
	return c.send(env)
}
