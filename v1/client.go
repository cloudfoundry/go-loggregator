// v1 provides a client to connect with the loggregtor v1 API
//
// Loggregator's v1 client library is better known to the Cloud Foundry
// community as Dropsonde (github.com/cloudfoundry/dropsonde). The code here
// wraps that library in the interest of consolidating all client code into
// a single library which includes both v1 and v2 clients.
package v1

import (
	"time"

	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"code.cloudfoundry.org/go-loggregator/v1/conversion"

	"github.com/cloudfoundry/dropsonde"
	"github.com/cloudfoundry/dropsonde/logs"
	"github.com/cloudfoundry/dropsonde/metrics"
	"github.com/cloudfoundry/sonde-go/events"
)

func NewClient() (*Client, error) {
	c := &Client{}

	return c, nil
}

type Client struct {
	tags map[string]*loggregator_v2.Value
}

func (c *Client) Send() error {
	return nil
}

func (c *Client) IncrementCounter(name string) error {
	return metrics.IncrementCounter(name)
}
func (c *Client) SendAppLog(appID, message, sourceType, sourceInstance string) error {
	return logs.SendAppLog(appID, message, sourceType, sourceInstance)
}

func (c *Client) SendAppErrorLog(appID, message, sourceType, sourceInstance string) error {
	return logs.SendAppErrorLog(appID, message, sourceType, sourceInstance)
}

func (c *Client) SendAppMetrics(m *events.ContainerMetric) error {
	return metrics.Send(m)
}

func (c *Client) SendDuration(name string, duration time.Duration) error {
	return c.SendComponentMetric(name, float64(duration), "nanos")
}

func (c *Client) SendMebiBytes(name string, mebibytes int) error {
	return c.SendComponentMetric(name, float64(mebibytes), "MiB")
}

func (c *Client) SendMetric(name string, value int) error {
	return c.SendComponentMetric(name, float64(value), "Metric")
}

func (c *Client) SendBytesPerSecond(name string, value float64) error {
	return c.SendComponentMetric(name, value, "B/s")
}

func (c *Client) SendRequestsPerSecond(name string, value float64) error {
	return c.SendComponentMetric(name, value, "Req/s")
}

func (c *Client) SendComponentMetric(name string, value float64, unit string) error {
	return metrics.SendValue(name, value, unit)
}

func (c *Client) EmitLog(message string, opts ...loggregator.EmitLogOption) {
	v2Envelope := &loggregator_v2.Envelope{
		Timestamp: time.Now().UnixNano(),
		Message: &loggregator_v2.Envelope_Log{
			Log: &loggregator_v2.Log{
				Payload: []byte(message),
				Type:    loggregator_v2.Log_ERR,
			},
		},
		Tags: make(map[string]*loggregator_v2.Value),
	}

	for _, o := range opts {
		o(v2Envelope)
	}
	c.emitEnvelopes(v2Envelope)
}

func (c *Client) EmitGauge(opts ...loggregator.EmitGaugeOption) {
	v2Envelope := &loggregator_v2.Envelope{
		Timestamp: time.Now().UnixNano(),
		Message: &loggregator_v2.Envelope_Gauge{
			Gauge: &loggregator_v2.Gauge{
				Metrics: make(map[string]*loggregator_v2.GaugeValue),
			},
		},
		Tags: make(map[string]*loggregator_v2.Value),
	}

	for _, o := range opts {
		o(v2Envelope)
	}
	c.emitEnvelopes(v2Envelope)
}

func (c *Client) EmitCounter(name string, opts ...loggregator.EmitCounterOption) {
	v2Envelope := &loggregator_v2.Envelope{
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

	for _, o := range opts {
		o(v2Envelope)
	}
	c.emitEnvelopes(v2Envelope)
}

func (c *Client) emitEnvelopes(v2Envelope *loggregator_v2.Envelope) {
	v1Envelopes := conversion.ToV1(v2Envelope)

	for _, e := range v1Envelopes {
		dropsonde.DefaultEmitter.EmitEnvelope(e)
	}
}
