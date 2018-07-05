package pulseemitter

import (
	"code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"github.com/golang/protobuf/proto"
	"time"
	"sync"
)

// TimerMetric is used by the pulse emitter to emit timer metrics to the
// LogClient.
type TimerMetric interface {
	// Record will add the start and stop time to the timer metric queue.
	Record(start, stop time.Time)

	// Emit sends the timer values to the LogClient.
	Emit(c LogClient)
}

// startStop contains a single timer record.
type startStop struct {
	start time.Time
	stop  time.Time
}

// timerMetric is used by the pulse emitter to emit timer metrics to the
// LogClient.
type timerMetric struct {
	name       string
	sourceID   string
	startStops []startStop
	tags       map[string]string

	mu *sync.Mutex
}

// NewTimerMetric returns a new timerMetric that keeps a queue of timer
// records to be emitted via a LogClient.
func NewTimerMetric(name, sourceID string, opts ...MetricOption) TimerMetric {
	t := &timerMetric{
		name:     name,
		sourceID: sourceID,
		tags:     make(map[string]string),

		mu: &sync.Mutex{},
	}

	for _, opt := range opts {
		opt(t.tags)
	}

	return t
}

// Record will add the start and stop time to the timer metric queue.
func (t *timerMetric) Record(start, stop time.Time) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.startStops = append(t.startStops, startStop{
		start: start,
		stop:  stop,
	})
}

// Emit will send the current tagging options and queue of metrics to
// the LogClient to be emitted.
func (t *timerMetric) Emit(c LogClient) {
	t.mu.Lock()
	waitingStartStops := t.startStops
	t.startStops = nil
	t.mu.Unlock()

	options := []loggregator.EmitTimerOption{
		t.sourceIDOption,
	}

	for k, v := range t.tags {
		options = append(options, loggregator.WithEnvelopeTag(k, v))
	}

	for _, startStop := range waitingStartStops {
		c.EmitTimer(t.name, startStop.start, startStop.stop, options...)
	}
}

func (t *timerMetric) sourceIDOption(p proto.Message) {
	env, ok := p.(*loggregator_v2.Envelope)
	if ok {
		env.SourceId = t.sourceID
	}
}
