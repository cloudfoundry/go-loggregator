package runtimeemitter

import (
	"runtime"
	"time"

	"code.cloudfoundry.org/go-loggregator"
)

// Sender is the interface of the client that can be used to emit gauge
// metrics.
type Sender interface {
	EmitGauge(opts ...loggregator.EmitGaugeOption)
}

// Emitter will emit a gauge with runtime stats via the sender on the given
// interval. default interval is 15 seconds.
type Emitter struct {
	interval time.Duration
	sender   Sender
}

// RuntimeEmitterOption is the option provides configuration for an Emitter.
type RuntimeEmitterOption func(e *Emitter)

// WithInterval returns a RuntimeEmitterOption to configure the interval at
// which the runtime emitter emits gauges.
func WithInterval(d time.Duration) RuntimeEmitterOption {
	return func(e *Emitter) {
		e.interval = d
	}
}

// New returns an Emitter that is configured with the given sender and
// RuntimeEmitterOptions.
func New(sender Sender, opts ...RuntimeEmitterOption) *Emitter {
	e := &Emitter{
		sender:   sender,
		interval: 10 * time.Second,
	}

	for _, o := range opts {
		o(e)
	}

	return e
}

// Run starts the ticker with the configured interval and emits a gauge on
// that interval. This method will block but the user may run in a go routine.
func (e *Emitter) Run() {
	for range time.Tick(e.interval) {
		memstats := &runtime.MemStats{}
		runtime.ReadMemStats(memstats)

		e.sender.EmitGauge(
			loggregator.WithGaugeValue("memoryStats.numBytesAllocatedHeap", float64(memstats.HeapAlloc), "Bytes"),
			loggregator.WithGaugeValue("memoryStats.numBytesAllocatedStack", float64(memstats.StackInuse), "Bytes"),
			loggregator.WithGaugeValue("memoryStats.lastGCPauseTimeNS", float64(memstats.PauseNs[(memstats.NumGC+255)%256]), "ns"),
			loggregator.WithGaugeValue("numGoRoutines", float64(runtime.NumGoroutine()), "Count"),
		)
	}
}
