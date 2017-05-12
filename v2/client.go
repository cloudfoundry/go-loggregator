// Package v2 provides a gRPC client to send data to the Loggregator v2 API.
package v2

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"code.cloudfoundry.org/go-loggregator/internal/loggregator_v2"
	"golang.org/x/net/context"
)

type Client struct {
	conn      loggregator_v2.IngressClient
	sender    loggregator_v2.Ingress_BatchSenderClient
	envelopes chan *loggregator_v2.Envelope

	batchMaxSize       uint
	batchFlushInterval time.Duration
	port               int

	logger Logger
}

type Option func(*Client)

func WithBatchMaxSize(maxSize uint) Option {
	return func(c *Client) {
		c.batchMaxSize = maxSize
	}
}

func WithBatchFlushInterval(d time.Duration) Option {
	return func(c *Client) {
		c.batchFlushInterval = d
	}
}

func WithPort(port int) Option {
	return func(c *Client) {
		c.port = port
	}
}

type Logger interface {
	Printf(string, ...interface{})
}

func WithLogger(l Logger) Option {
	return func(c *Client) {
		c.logger = l
	}
}

func NewClient(tlsConfig *tls.Config, opts ...Option) (*Client, error) {
	client := &Client{
		envelopes:          make(chan *loggregator_v2.Envelope, 100),
		batchMaxSize:       100,
		batchFlushInterval: time.Second,
		port:               3458,
		logger:             log.New(ioutil.Discard, "", 0),
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

type EmitLogOption func(*loggregator_v2.Envelope)

func WithAppInfo(appID, sourceType, sourceInstance string) EmitLogOption {
	return func(e *loggregator_v2.Envelope) {
		e.SourceId = appID
		e.InstanceId = sourceInstance

		// TODO: don't blow away the tags
		e.Tags = map[string]*loggregator_v2.Value{
			"source_type":     &loggregator_v2.Value{Data: &loggregator_v2.Value_Text{Text: sourceType}},
			"source_instance": &loggregator_v2.Value{Data: &loggregator_v2.Value_Text{Text: sourceInstance}},
		}
	}
}

func WithStdout() EmitLogOption {
	return func(e *loggregator_v2.Envelope) {
		e.GetLog().Type = loggregator_v2.Log_OUT
	}
}

func (c *Client) EmitLog(message string, opts ...EmitLogOption) {
	e := &loggregator_v2.Envelope{
		Timestamp: time.Now().UnixNano(),
		Message: &loggregator_v2.Envelope_Log{
			Log: &loggregator_v2.Log{
				Payload: []byte(message),
				Type:    loggregator_v2.Log_ERR,
			},
		},
	}

	for _, o := range opts {
		o(e)
	}

	c.envelopes <- e
}

type EmitGaugeOption func(*loggregator_v2.Envelope)

func WithGaugeAppInfo(appID string) EmitGaugeOption {
	return func(e *loggregator_v2.Envelope) {
		e.SourceId = appID
	}
}

func WithGaugeValue(name string, value float64, unit string) EmitGaugeOption {
	return func(e *loggregator_v2.Envelope) {
		e.GetGauge().Metrics[name] = &loggregator_v2.GaugeValue{Value: value, Unit: unit}
	}
}

func (c *Client) EmitGauge(opts ...EmitGaugeOption) {
	e := &loggregator_v2.Envelope{
		Timestamp: time.Now().UnixNano(),
		Message: &loggregator_v2.Envelope_Gauge{
			Gauge: &loggregator_v2.Gauge{
				Metrics: make(map[string]*loggregator_v2.GaugeValue),
			},
		},
	}

	for _, o := range opts {
		o(e)
	}

	c.envelopes <- e
}

func (c *Client) EmitCounter(name string) {
	e := &loggregator_v2.Envelope{
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

	c.envelopes <- e
}

func (c *Client) startSender() {
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

func (c *Client) flush(batch []*loggregator_v2.Envelope) {
	if c.sender == nil {
		var err error
		c.sender, err = c.conn.BatchSender(context.TODO())
		if err != nil {
			c.logger.Printf("Error while flushing: %s", err)
			return
		}
	}

	err := c.sender.Send(&loggregator_v2.EnvelopeBatch{Batch: batch})
	if err != nil {
		c.logger.Printf("Error while flushing: %s", err)
		c.sender = nil
		return
	}

	return
}
