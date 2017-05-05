package loggregator_v2

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/cloudfoundry/sonde-go/events"
)

type envelopeWithResponseChannel struct {
	envelope *Envelope
	errCh    chan error
}

func newGrpcClient(config MetronConfig) (*grpcClient, error) {
	tlsConfig, err := newTLSConfig(
		config.CACertPath,
		config.CertPath,
		config.KeyPath,
	)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%d", config.APIPort),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	)
	if err != nil {
		return nil, err
	}

	client := &grpcClient{
		ingressClient: NewIngressClient(conn),
		config:        config,
		envelopes:     make(chan *envelopeWithResponseChannel),
	}

	go client.startSender()

	return client, nil
}

type grpcClient struct {
	ingressClient IngressClient
	sender        Ingress_BatchSenderClient
	envelopes     chan *envelopeWithResponseChannel
	config        MetronConfig
}

func (c *grpcClient) startSender() {
	for {
		envelopeWithResponseChannel := <-c.envelopes
		envelope := envelopeWithResponseChannel.envelope
		errCh := envelopeWithResponseChannel.errCh
		if c.sender == nil {
			var err error
			c.sender, err = c.ingressClient.BatchSender(context.Background())
			if err != nil {
				errCh <- err
				continue
			}
		}
		err := c.sender.Send(&EnvelopeBatch{Batch: []*Envelope{envelope}})
		if err != nil {
			c.sender = nil
		}
		errCh <- err
	}
}

func (c *grpcClient) send(envelope *Envelope) error {
	if envelope.Tags == nil {
		envelope.Tags = make(map[string]*Value)
	}
	envelope.Tags["deployment"] = newTextValue(c.config.JobDeployment)
	envelope.Tags["job"] = newTextValue(c.config.JobName)
	envelope.Tags["index"] = newTextValue(c.config.JobIndex)
	envelope.Tags["ip"] = newTextValue(c.config.JobIP)
	envelope.Tags["origin"] = newTextValue(c.config.JobOrigin)

	e := &envelopeWithResponseChannel{
		envelope: envelope,
		errCh:    make(chan error),
	}
	defer close(e.errCh)

	c.envelopes <- e
	err := <-e.errCh
	return err
}

func (c *grpcClient) SendAppLog(appID, message, sourceType, sourceInstance string) error {
	return c.send(createLogEnvelope(appID, message, sourceType, sourceInstance, Log_OUT))
}

func (c *grpcClient) SendAppErrorLog(appID, message, sourceType, sourceInstance string) error {
	return c.send(createLogEnvelope(appID, message, sourceType, sourceInstance, Log_ERR))
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

func (c *grpcClient) sendGauge(metrics map[string]*GaugeValue) error {
	return c.send(&Envelope{
		Timestamp: time.Now().UnixNano(),
		Message: &Envelope_Gauge{
			Gauge: &Gauge{
				Metrics: metrics,
			},
		},
	})
}

func (c *grpcClient) SendDuration(name string, duration time.Duration) error {
	metrics := make(map[string]*GaugeValue)
	metrics[name] = &GaugeValue{
		Unit:  "nanos",
		Value: float64(duration),
	}
	return c.sendGauge(metrics)
}

func (c *grpcClient) SendMebiBytes(name string, mebibytes int) error {
	metrics := make(map[string]*GaugeValue)
	metrics[name] = &GaugeValue{
		Unit:  "MiB",
		Value: float64(mebibytes),
	}
	return c.sendGauge(metrics)
}

func (c *grpcClient) SendMetric(name string, value int) error {
	metrics := make(map[string]*GaugeValue)
	metrics[name] = &GaugeValue{
		Unit:  "Metric",
		Value: float64(value),
	}
	return c.sendGauge(metrics)
}

func (c *grpcClient) SendBytesPerSecond(name string, value float64) error {
	metrics := make(map[string]*GaugeValue)
	metrics[name] = &GaugeValue{
		Unit:  "B/s",
		Value: float64(value),
	}
	return c.sendGauge(metrics)
}

func (c *grpcClient) SendRequestsPerSecond(name string, value float64) error {
	metrics := make(map[string]*GaugeValue)
	metrics[name] = &GaugeValue{
		Unit:  "Req/s",
		Value: float64(value),
	}
	return c.sendGauge(metrics)
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
