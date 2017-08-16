// v1 provides a client to connect with the loggregtor v1 API
//
// Loggregator's v1 client library is better known to the Cloud Foundry
// community as Dropsonde (github.com/cloudfoundry/dropsonde). The code here
// wraps that library in the interest of consolidating all client code into
// a single library which includes both v1 and v2 clients.
package v1

import (
	"io/ioutil"
	"log"
	"time"

	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"code.cloudfoundry.org/go-loggregator/v1/conversion"

	"github.com/cloudfoundry/dropsonde"
)

type ClientOption func(*Client)

// WithTag allows for the configuration of arbitrary string value
// metadata which will be included in all data sent to Loggregator
func WithTag(name, value string) ClientOption {
	return func(c *Client) {
		c.tags[name] = value
	}
}

// WithLogger allows for the configuration of a logger.
// By default, the logger is disabled.
func WithLogger(l loggregator.Logger) ClientOption {
	return func(c *Client) {
		c.logger = l
	}
}

// Client represents an emitter into loggregator. It should be created with
// the NewClient constructor.
type Client struct {
	tags   map[string]string
	logger loggregator.Logger
}

// NewClient creates a v1 loggregator client. This is a wrapper around the
// dropsonde package that will write envelopes to loggregator over UDP. Before
// calling NewClient you should call dropsonde.Initialize.
func NewClient(opts ...ClientOption) (*Client, error) {
	c := &Client{
		tags:   make(map[string]string),
		logger: log.New(ioutil.Discard, "", 0),
	}

	for _, o := range opts {
		o(c)
	}

	return c, nil
}

// EmitLog sends a message to loggregator.
func (c *Client) EmitLog(message string, opts ...loggregator.EmitLogOption) {
	v2Envelope := &loggregator_v2.Envelope{
		Timestamp: time.Now().UnixNano(),
		Message: &loggregator_v2.Envelope_Log{
			Log: &loggregator_v2.Log{
				Payload: []byte(message),
				Type:    loggregator_v2.Log_ERR,
			},
		},
		Tags: make(map[string]string),
	}

	for _, o := range opts {
		o(v2Envelope)
	}
	c.emitEnvelopes(v2Envelope)
}

// EmitGauge sends the configured gauge values to loggregator.
// If no EmitGaugeOption values are present, no envelopes will be emitted.
func (c *Client) EmitGauge(opts ...loggregator.EmitGaugeOption) {
	v2Envelope := &loggregator_v2.Envelope{
		Timestamp: time.Now().UnixNano(),
		Message: &loggregator_v2.Envelope_Gauge{
			Gauge: &loggregator_v2.Gauge{
				Metrics: make(map[string]*loggregator_v2.GaugeValue),
			},
		},
		Tags: make(map[string]string),
	}

	for _, o := range opts {
		o(v2Envelope)
	}
	c.emitEnvelopes(v2Envelope)
}

// EmitCounter sends a counter envelope with a delta of 1.
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
		Tags: make(map[string]string),
	}

	for _, o := range opts {
		o(v2Envelope)
	}
	c.emitEnvelopes(v2Envelope)
}

func (c *Client) emitEnvelopes(v2Envelope *loggregator_v2.Envelope) {
	for k, v := range c.tags {
		v2Envelope.Tags[k] = v
	}
	v2Envelope.Tags["origin"] = dropsonde.DefaultEmitter.Origin()

	for _, e := range conversion.ToV1(v2Envelope) {
		err := dropsonde.DefaultEmitter.EmitEnvelope(e)
		if err != nil {
			c.logger.Printf("Failed to emit envelope: %s", err)
		}
	}
}
