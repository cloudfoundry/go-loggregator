package pulseemitter_test

import (
	"sync"

	"code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/pulseemitter"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("CounterMetric", func() {
	Context("Emit", func() {
		It("prepares an envelope for delivery", func() {
			metric := pulseemitter.NewCounterMetric(
				"name",
				"my-source-id",
				pulseemitter.WithVersion(1, 2),
			)

			metric.Increment(10)

			spy := newSpyLogClient()
			metric.Emit(spy)
			Expect(spy.CounterName()).To(Equal("name"))

			e := &loggregator_v2.Envelope{
				Message: &loggregator_v2.Envelope_Counter{
					Counter: &loggregator_v2.Counter{},
				},
				Tags: make(map[string]string),
			}
			for _, o := range spy.CounterOpts() {
				o(e)
			}

			Expect(e.GetCounter().GetDelta()).To(Equal(uint64(10)))
			Expect(e.Tags["metric_version"]).To(Equal("1.2"))
		})

		It("decrements its value on success", func() {
			metric := pulseemitter.NewCounterMetric("name", "my-source-id")
			spy := newSpyLogClient()

			metric.Increment(10)
			metric.Emit(spy)

			metric.Emit(spy)
			e := &loggregator_v2.Envelope{
				Message: &loggregator_v2.Envelope_Counter{
					Counter: &loggregator_v2.Counter{},
				},
			}

			for _, o := range spy.counterOpts {
				o(e)
			}

			Expect(e.GetCounter().GetDelta()).To(Equal(uint64(0)))
		})
	})
})

type spyTimer struct {
	name  string
	start time.Time
	stop  time.Time
	opts  []loggregator.EmitTimerOption
}

type spyLogClient struct {
	mu             sync.Mutex
	counterName    string
	counterOpts    []loggregator.EmitCounterOption
	gaugeOpts      []loggregator.EmitGaugeOption
	gaugeCallCount int
	timers         []spyTimer
}

func newSpyLogClient() *spyLogClient {
	return &spyLogClient{}
}

func (s *spyLogClient) EmitCounter(name string, opts ...loggregator.EmitCounterOption) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counterName = name
	s.counterOpts = opts
}

func (s *spyLogClient) EmitGauge(opts ...loggregator.EmitGaugeOption) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gaugeCallCount++
	s.gaugeOpts = opts
}

func (s *spyLogClient) EmitTimer(name string, start, stop time.Time, opts ...loggregator.EmitTimerOption) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.timers = append(s.timers, spyTimer{
		name:  name,
		start: start,
		stop:  stop,
		opts:  opts,
	})
}

func (s *spyLogClient) CounterName() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.counterName
}

func (s *spyLogClient) CounterOpts() []loggregator.EmitCounterOption {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.counterOpts
}

func (s *spyLogClient) GaugeOpts() []loggregator.EmitGaugeOption {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.gaugeOpts
}

func (s *spyLogClient) GaugeCallCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.gaugeCallCount
}

func (s *spyLogClient) ResetTimers() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.timers = nil
}

func (s *spyLogClient) Timers() []spyTimer {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.timers
}
