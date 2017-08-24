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

	"github.com/cloudfoundry/dropsonde"
	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gogo/protobuf/proto"
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
	w := loggregator.EnvelopeWrapper{
		Messages: []*events.Envelope{
			{
				Timestamp: proto.Int64(time.Now().UnixNano()),
				EventType: events.Envelope_LogMessage.Enum(),
				LogMessage: &events.LogMessage{
					MessageType: events.LogMessage_ERR.Enum(),
					Message:     []byte(message),
					Timestamp:   proto.Int64(time.Now().UnixNano()),
				},
			},
		},
		Tags: make(map[string]string),
	}

	for _, o := range opts {
		o(&w)
	}
	w.Messages[0].Tags = w.Tags

	c.emitEnvelope(w)
}

// EmitGauge sends the configured gauge values to loggregator.
// If no EmitGaugeOption values are present, no envelopes will be emitted.
func (c *Client) EmitGauge(opts ...loggregator.EmitGaugeOption) {
	w := loggregator.EnvelopeWrapper{
		Tags: make(map[string]string),
	}

	for _, o := range opts {
		o(&w)
	}

	for _, e := range w.Messages {
		e.Timestamp = proto.Int64(time.Now().UnixNano())
		e.EventType = events.Envelope_ValueMetric.Enum()
		e.Tags = w.Tags
	}

	c.emitEnvelope(w)
}

// EmitCounter sends a counter envelope with a delta of 1.
func (c *Client) EmitCounter(name string, opts ...loggregator.EmitCounterOption) {
	w := loggregator.EnvelopeWrapper{
		Messages: []*events.Envelope{
			{
				Timestamp: proto.Int64(time.Now().UnixNano()),
				EventType: events.Envelope_CounterEvent.Enum(),
				CounterEvent: &events.CounterEvent{
					Name:  proto.String(name),
					Delta: proto.Uint64(1),
				},
			},
		},
		Tags: make(map[string]string),
	}

	for _, o := range opts {
		o(&w)
	}

	w.Messages[0].Tags = w.Tags

	c.emitEnvelope(w)
}

func (c *Client) emitEnvelope(w loggregator.EnvelopeWrapper) {
	for _, e := range w.Messages {
		e.Origin = proto.String(dropsonde.DefaultEmitter.Origin())
		for k, v := range c.tags {
			e.Tags[k] = v
		}

		err := dropsonde.DefaultEmitter.EmitEnvelope(e)
		if err != nil {
			c.logger.Printf("Failed to emit envelope: %s", err)
		}
	}
}
