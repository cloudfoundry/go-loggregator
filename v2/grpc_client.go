package v2

import (
	"crypto/tls"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"code.cloudfoundry.org/go-loggregator/internal/loggregator_v2"
	"github.com/cloudfoundry/sonde-go/events"
	"golang.org/x/net/context"
)

type grpcClient struct {
	conn      loggregator_v2.IngressClient
	sender    loggregator_v2.Ingress_BatchSenderClient
	envelopes chan *loggregator_v2.Envelope

	batchMaxSize       uint
	batchFlushInterval time.Duration
	port               int
}

type Option func(*grpcClient)

func WithBatchMaxSize(maxSize uint) Option {
	return func(c *grpcClient) {
		c.batchMaxSize = maxSize
	}
}

func WithBatchFlushInterval(d time.Duration) Option {
	return func(c *grpcClient) {
		c.batchFlushInterval = d
	}
}

func WithPort(port int) Option {
	return func(c *grpcClient) {
		c.port = port
	}
}

func NewClient(tlsConfig *tls.Config, opts ...Option) (*grpcClient, error) {
	client := &grpcClient{
		envelopes:          make(chan *loggregator_v2.Envelope),
		batchMaxSize:       100,
		batchFlushInterval: time.Second,
		port:               8082,
	}

	for _, o := range opts {
		o(client)
	}

	conn, err := grpc.Dial(
		fmt.Sprintf("localhost:%d", client.port),
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	)
	if err != nil {
		return nil, err
	}
	client.conn = loggregator_v2.NewIngressClient(conn)

	go client.startSender()

	return client, nil
}

func (c *grpcClient) SendAppLog(appID, message, sourceType, sourceInstance string) {
	c.send(createLogEnvelope(appID, message, sourceType, sourceInstance, loggregator_v2.Log_OUT))
}

func (c *grpcClient) SendAppErrorLog(appID, message, sourceType, sourceInstance string) {
	c.send(createLogEnvelope(appID, message, sourceType, sourceInstance, loggregator_v2.Log_ERR))
}

func (c *grpcClient) SendAppMetrics(m *events.ContainerMetric) {
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
	c.send(env)
}

func (c *grpcClient) SendDuration(name string, duration time.Duration) {
	metrics := make(map[string]*loggregator_v2.GaugeValue)
	metrics[name] = &loggregator_v2.GaugeValue{
		Unit:  "nanos",
		Value: float64(duration),
	}
	c.sendGauge(metrics)
}

func (c *grpcClient) SendMebiBytes(name string, mebibytes int) {
	metrics := make(map[string]*loggregator_v2.GaugeValue)
	metrics[name] = &loggregator_v2.GaugeValue{
		Unit:  "MiB",
		Value: float64(mebibytes),
	}
	c.sendGauge(metrics)
}

func (c *grpcClient) SendMetric(name string, value int) {
	metrics := make(map[string]*loggregator_v2.GaugeValue)
	metrics[name] = &loggregator_v2.GaugeValue{
		Unit:  "Metric",
		Value: float64(value),
	}
	c.sendGauge(metrics)
}

func (c *grpcClient) SendBytesPerSecond(name string, value float64) {
	metrics := make(map[string]*loggregator_v2.GaugeValue)
	metrics[name] = &loggregator_v2.GaugeValue{
		Unit:  "B/s",
		Value: float64(value),
	}
	c.sendGauge(metrics)
}

func (c *grpcClient) SendRequestsPerSecond(name string, value float64) {
	metrics := make(map[string]*loggregator_v2.GaugeValue)
	metrics[name] = &loggregator_v2.GaugeValue{
		Unit:  "Req/s",
		Value: float64(value),
	}
	c.sendGauge(metrics)
}

func (c *grpcClient) IncrementCounter(name string) {
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

	c.send(env)
}

func (c *grpcClient) startSender() {
	t := time.NewTimer(c.batchFlushInterval)

	var batch []*loggregator_v2.Envelope
	for {
		select {
		case env := <-c.envelopes:
			batch = append(batch, env)

			if len(batch) >= int(c.batchMaxSize) {
				c.flush(batch)
				batch = nil
			}

			if !t.Stop() {
				<-t.C
			}
		case <-t.C:
			if len(batch) > 0 {
				c.flush(batch)
				batch = nil
			}
		}
		t.Reset(c.batchFlushInterval)
	}
}

func (c *grpcClient) flush(batch []*loggregator_v2.Envelope) error {
	if c.sender == nil {
		var err error
		c.sender, err = c.conn.BatchSender(context.TODO())
		if err != nil {
			return err
		}
	}

	err := c.sender.Send(&loggregator_v2.EnvelopeBatch{Batch: batch})
	if err != nil {
		c.sender = nil
		return err
	}

	return nil
}

func (c *grpcClient) send(envelope *loggregator_v2.Envelope) {
	if envelope.Tags == nil {
		envelope.Tags = make(map[string]*loggregator_v2.Value)
	}
	c.envelopes <- envelope
}

func (c *grpcClient) sendGauge(metrics map[string]*loggregator_v2.GaugeValue) {
	c.send(&loggregator_v2.Envelope{
		Timestamp: time.Now().UnixNano(),
		Message: &loggregator_v2.Envelope_Gauge{
			Gauge: &loggregator_v2.Gauge{
				Metrics: metrics,
			},
		},
	})
}
