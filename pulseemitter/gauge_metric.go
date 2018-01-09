package pulseemitter

import (
	"sync/atomic"

	loggregator "code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/golang/protobuf/proto"
)

// GaugeMetric is used by the pulse emitter to emit gauge metrics to the
// LoggClient.
type GaugeMetric interface {
	// Set sets the current value of the gauge metric.
	Set(n int64)

	// Emit sends the counter values to the LoggClient.
	Emit(c LoggClient)
}

// gaugeMetric is used by the pulse emitter to emit gauge metrics to the
// LoggClient.
type gaugeMetric struct {
	name     string
	unit     string
	sourceID string
	value    int64
	tags     map[string]string
}

// NewGaugeMetric returns a new gaugeMetric that has a value that can be set
// and emitted via a LoggClient.
func NewGaugeMetric(name, unit, sourceID string, opts ...MetricOption) GaugeMetric {
	g := &gaugeMetric{
		name:     name,
		unit:     unit,
		sourceID: sourceID,
		tags:     make(map[string]string),
	}

	for _, opt := range opts {
		opt(g.tags)
	}

	return g
}

// Set will set the current value of the gauge metric to the given number.
func (g *gaugeMetric) Set(n int64) {
	atomic.SwapInt64(&g.value, n)
}

// Emit will send the current value and tagging options to the LoggClient to
// be emitted.
func (g *gaugeMetric) Emit(c LoggClient) {
	options := []loggregator.EmitGaugeOption{
		loggregator.WithGaugeValue(
			g.name,
			float64(atomic.LoadInt64(&g.value)),
			g.unit,
		),
		g.sourceIDOption,
	}

	for k, v := range g.tags {
		options = append(options, loggregator.WithEnvelopeTag(k, v))
	}

	c.EmitGauge(options...)
}

func (g *gaugeMetric) sourceIDOption(p proto.Message) {
	env, ok := p.(*loggregator_v2.Envelope)
	if ok {
		env.SourceId = g.sourceID
	}
}
