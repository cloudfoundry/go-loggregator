package v2

import (
	"time"

	"code.cloudfoundry.org/go-loggregator/internal/loggregator_v2"
	"github.com/cloudfoundry/sonde-go/events"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type envelopeWithResponseChannel struct {
	envelope *loggregator_v2.Envelope
	errCh    chan error
}

type grpcClient struct {
	batchStreamer BatchStreamer
	sender        loggregator_v2.Ingress_BatchSenderClient
	envelopes     chan *envelopeWithResponseChannel
	jobOpts       JobOpts
}

type JobOpts struct {
	Deployment string
	Name       string
	Index      string
	IP         string
	Origin     string
}

type BatchStreamer interface {
	BatchSender(ctx context.Context, opts ...grpc.CallOption) (loggregator_v2.Ingress_BatchSenderClient, error)
}

type v2Opt func(*grpcClient)

func WithJobOpts(j JobOpts) func(*grpcClient) {
	return func(c *grpcClient) {
		c.jobOpts = j
	}
}

func NewClient(b BatchStreamer, opts ...v2Opt) (*grpcClient, error) {
	client := &grpcClient{
		batchStreamer: b,
		envelopes:     make(chan *envelopeWithResponseChannel),
	}

	for _, o := range opts {
		o(client)
	}

	go client.startSender()

	return client, nil
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
			c.sender, err = c.batchStreamer.BatchSender(context.TODO())
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
	envelope.Tags["deployment"] = newTextValue(c.jobOpts.Deployment)
	envelope.Tags["job"] = newTextValue(c.jobOpts.Name)
	envelope.Tags["index"] = newTextValue(c.jobOpts.Index)
	envelope.Tags["ip"] = newTextValue(c.jobOpts.IP)
	envelope.Tags["origin"] = newTextValue(c.jobOpts.Origin)

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
