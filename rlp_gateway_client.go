package loggregator

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/gogo/protobuf/jsonpb"
	"golang.org/x/net/context"
)

type RLPGatewayClient struct {
	addr string
	log  *log.Logger
	doer Doer
}

func NewRLPGatewayClient(addr string, opts ...RLPGatewayClientOption) *RLPGatewayClient {
	c := &RLPGatewayClient{
		addr: addr,
		log:  log.New(ioutil.Discard, "", 0),
		doer: http.DefaultClient,
	}

	for _, o := range opts {
		o(c)
	}

	return c
}

// RLPGatewayClientOption is the type of a configurable client option.
type RLPGatewayClientOption func(*RLPGatewayClient)

// WithRLPGatewayClientLogger returns a RLPGatewayClientOption to configure
// the logger of the RLPGatewayClient. It defaults to a silent logger.
func WithRLPGatewayClientLogger(log *log.Logger) RLPGatewayClientOption {
	return func(c *RLPGatewayClient) {
		c.log = log
	}
}

// WithRLPGatewayClientLogger returns a RLPGatewayClientOption to configure
// the HTTP client. It defaults to the http.DefaultClient.
func WithRLPGatewayHTTPClient(d Doer) RLPGatewayClientOption {
	return func(c *RLPGatewayClient) {
		c.doer = d
	}
}

// Doer is used to make HTTP requests to the RLP Gateway.
type Doer interface {
	// Do is a implementation of the http.Client's Do method.
	Do(*http.Request) (*http.Response, error)
}

// Stream returns a new EnvelopeStream for the given context and request. The
// lifecycle of the EnvelopeStream is managed by the given context. If the
// underlying SSE stream dies, it attempts to reconnect until the context
// is done. Any errors are logged via the client's logger.
func (c *RLPGatewayClient) Stream(ctx context.Context, req *loggregator_v2.EgressBatchRequest) EnvelopeStream {
	es := make(chan *loggregator_v2.Envelope, 100)
	go func() {
		defer close(es)
		for ctx.Err() == nil {
			c.connect(ctx, es, req)
		}
	}()

	return func() []*loggregator_v2.Envelope {
		var batch []*loggregator_v2.Envelope
		for {
			select {
			case <-ctx.Done():
				return nil
			case e, ok := <-es:
				if !ok {
					return nil
				}
				batch = append(batch, e)
			default:
				if len(batch) > 0 {
					return batch
				}
			}
		}
	}
}

func (c *RLPGatewayClient) connect(ctx context.Context, es chan<- *loggregator_v2.Envelope, logReq *loggregator_v2.EgressBatchRequest) {
	readAddr := fmt.Sprintf("%s/v2/read%s", c.addr, c.buildQuery(logReq))
	req, err := http.NewRequest(http.MethodGet, readAddr, nil)
	if err != nil {
		c.log.Fatalf("failed to build request %s", err)
	}
	req.Header.Set("Content-Type", "text/event-stream")

	resp, err := c.doer.Do(req)
	if err != nil {
		c.log.Printf("error making request: %s", err)
		return
	}

	defer func() {
		io.Copy(ioutil.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.log.Printf("failed to read body: %s", err)
			return
		}
		c.log.Printf("unexpected status code %d: %s", resp.StatusCode, body)
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data:") {
			continue
		}
		line = line[len("data:"):]

		var e loggregator_v2.Envelope
		if err := jsonpb.Unmarshal(strings.NewReader(line), &e); err != nil {
			c.log.Printf("failed to unmarshal envelope: %s", err)
			continue
		}

		select {
		case <-ctx.Done():
			return
		case es <- &e:
		}
	}

	if scanner.Err() != nil {
		c.log.Printf("failed while reading stream: %s", scanner.Err())
	}
}

func (c *RLPGatewayClient) buildQuery(req *loggregator_v2.EgressBatchRequest) string {
	var query []string
	if req.GetShardId() != "" {
		query = append(query, "shard_id="+req.GetShardId())
	}

	if req.GetDeterministicName() != "" {
		query = append(query, "deterministic_name="+req.GetDeterministicName())
	}

	for _, selector := range req.GetSelectors() {
		if selector.GetSourceId() != "" {
			query = append(query, "source_id="+selector.GetSourceId())
		}

		switch selector.Message.(type) {
		case *loggregator_v2.Selector_Log:
			query = append(query, "log")
		case *loggregator_v2.Selector_Counter:
			if selector.GetCounter().GetName() != "" {
				query = append(query, "counter.name="+selector.GetCounter().GetName())
				continue
			}
			query = append(query, "counter")
		case *loggregator_v2.Selector_Gauge:
			if len(selector.GetGauge().GetNames()) > 1 {
				// TODO: This is a mistake in the gateway.
				panic("This is not yet supported")
			}

			if len(selector.GetGauge().GetNames()) != 0 {
				query = append(query, "gauge.name="+selector.GetGauge().GetNames()[0])
				continue
			}
			query = append(query, "gauge")
		case *loggregator_v2.Selector_Timer:
			query = append(query, "timer")
		case *loggregator_v2.Selector_Event:
			query = append(query, "event")
		}
	}

	namedCounter := containsPrefix(query, "counter.name")
	namedGauge := containsPrefix(query, "gauge.name")

	if namedCounter {
		query = filter(query, "counter")
	}

	if namedGauge {
		query = filter(query, "gauge")
	}

	if len(query) == 0 {
		return ""
	}

	return "?" + strings.Join(query, "&")
}

func containsPrefix(arr []string, prefix string) bool {
	for _, i := range arr {
		if strings.HasPrefix(i, prefix) {
			return true
		}
	}
	return false
}

func filter(arr []string, target string) []string {
	var filtered []string
	for _, i := range arr {
		if i != target {
			filtered = append(filtered, i)
		}
	}
	return filtered
}
