package v2

import (
	"context"
	"fmt"
	"time"

	"code.cloudfoundry.org/go-loggregator/loggregator_v2"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/cloudfoundry/sonde-go/events"
)

func NewGrpcClient(config MetronConfig) (*grpcClient, error) {
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
		ingressClient: loggregator_v2.NewIngressClient(conn),
		config:        config,
		envelopes:     make(chan *envelopeWithResponseChannel),
	}

	go client.startSender()

	return client, nil
}

type MetronConfig struct {
	UseV2API      bool   `json:"loggregator_use_v2_api"`
	APIPort       int    `json:"loggregator_api_port"`
	CACertPath    string `json:"loggregator_ca_path"`
	CertPath      string `json:"loggregator_cert_path"`
	KeyPath       string `json:"loggregator_key_path"`
	JobDeployment string `json:"loggregator_job_deployment"`
	JobName       string `json:"loggregator_job_name"`
	JobIndex      string `json:"loggregator_job_index"`
	JobIP         string `json:"loggregator_job_ip"`
	JobOrigin     string `json:"loggregator_job_origin"`
	DropsondePort int    `json:"dropsonde_port"`
}

type envelopeWithResponseChannel struct {
	envelope *loggregator_v2.Envelope
	errCh    chan error
}

type grpcClient struct {
	ingressClient loggregator_v2.IngressClient
	sender        loggregator_v2.Ingress_BatchSenderClient
	envelopes     chan *envelopeWithResponseChannel
	config        MetronConfig
}

func (c *grpcClient) SendAppLog(appID, message, sourceType, sourceInstance string) error {
	return c.send(createLogEnvelope(appID, message, sourceType, sourceInstance, loggregator_v2.Log_OUT))
}

func (c *grpcClient) SendAppErrorLog(appID, message, sourceType, sourceInstance string) error {
	return c.send(createLogEnvelope(appID, message, sourceType, sourceInstance, loggregator_v2.Log_ERR))
}

func (c *grpcClient) SendAppMetrics(m *events.ContainerMetric) error {
	env := &loggregator_v2.Envelope{
		Timestamp: time.Now().UnixNano(),
		SourceId:  m.GetApplicationId(),
		Message: &loggregator_v2.Envelope_Gauge{
			Gauge: &loggregator_v2.Gauge{
				Metrics: map[string]*loggregator_v2.GaugeValue{
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
	metrics := make(map[string]*loggregator_v2.GaugeValue)
	metrics[name] = &loggregator_v2.GaugeValue{
		Unit:  "nanos",
		Value: float64(duration),
	}
	return c.sendGauge(metrics)
}

func (c *grpcClient) SendMebiBytes(name string, mebibytes int) error {
	metrics := make(map[string]*loggregator_v2.GaugeValue)
	metrics[name] = &loggregator_v2.GaugeValue{
		Unit:  "MiB",
		Value: float64(mebibytes),
	}
	return c.sendGauge(metrics)
}

func (c *grpcClient) SendMetric(name string, value int) error {
	metrics := make(map[string]*loggregator_v2.GaugeValue)
	metrics[name] = &loggregator_v2.GaugeValue{
		Unit:  "Metric",
		Value: float64(value),
	}
	return c.sendGauge(metrics)
}

func (c *grpcClient) SendBytesPerSecond(name string, value float64) error {
	metrics := make(map[string]*loggregator_v2.GaugeValue)
	metrics[name] = &loggregator_v2.GaugeValue{
		Unit:  "B/s",
		Value: float64(value),
	}
	return c.sendGauge(metrics)
}

func (c *grpcClient) SendRequestsPerSecond(name string, value float64) error {
	metrics := make(map[string]*loggregator_v2.GaugeValue)
	metrics[name] = &loggregator_v2.GaugeValue{
		Unit:  "Req/s",
		Value: float64(value),
	}
	return c.sendGauge(metrics)
}

func (c *grpcClient) IncrementCounter(name string) error {
	env := &loggregator_v2.Envelope{
		Timestamp: time.Now().UnixNano(),
		Message: &loggregator_v2.Envelope_Counter{
			Counter: &loggregator_v2.Counter{
				Name: name,
				Value: &loggregator_v2.Counter_Delta{
					Delta: uint64(1),
				},
			},
		},
	}
	return c.send(env)
}

func (c *grpcClient) startSender() {
	for {
		envelopeWithResponseChannel := <-c.envelopes
		envelope := envelopeWithResponseChannel.envelope
		errCh := envelopeWithResponseChannel.errCh
		if c.sender == nil {
			var err error
			c.sender, err = c.ingressClient.BatchSender(context.TODO())
			if err != nil {
				errCh <- err
				continue
			}
		}
		err := c.sender.Send(&loggregator_v2.EnvelopeBatch{Batch: []*loggregator_v2.Envelope{envelope}})
		if err != nil {
			c.sender = nil
		}
		errCh <- err
	}
}

func (c *grpcClient) send(envelope *loggregator_v2.Envelope) error {
	if envelope.Tags == nil {
		envelope.Tags = make(map[string]*loggregator_v2.Value)
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

func (c *grpcClient) sendGauge(metrics map[string]*loggregator_v2.GaugeValue) error {
	return c.send(&loggregator_v2.Envelope{
		Timestamp: time.Now().UnixNano(),
		Message: &loggregator_v2.Envelope_Gauge{
			Gauge: &loggregator_v2.Gauge{
				Metrics: metrics,
			},
		},
	})
}
