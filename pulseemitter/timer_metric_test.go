package pulseemitter_test

import (
	"code.cloudfoundry.org/go-loggregator/pulseemitter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"sync"
)

var _ = Describe("TimerMetric", func() {
	It("prepares the envelope for delivery", func() {
		t := pulseemitter.NewTimerMetric(
			"some-timer",
			"my-source-id",
			pulseemitter.WithVersion(1, 2),
		)

		startTime := time.Now().Add(-time.Minute)
		stopTime := time.Now()
		t.Record(startTime, stopTime)

		spy := newSpyLogClient()
		t.Emit(spy)

		timers := spy.Timers()
		Expect(timers).To(HaveLen(1))
		Expect(timers[0].name).To(Equal("some-timer"))
		Expect(timers[0].start).To(Equal(startTime))
		Expect(timers[0].stop).To(Equal(stopTime))

		e := &loggregator_v2.Envelope{
			Message: &loggregator_v2.Envelope_Timer{
				Timer: &loggregator_v2.Timer{},
			},
			Tags: make(map[string]string),
		}
		for _, o := range timers[0].opts {
			o(e)
		}
		Expect(e.GetTags()).To(HaveKey("metric_version"))
		Expect(e.GetTags()["metric_version"]).To(Equal("1.2"))
	})

	It("emits all new timers", func() {
		t := pulseemitter.NewTimerMetric(
			"some-timer",
			"my-source-id",
		)
		t.Record(time.Now(), time.Now())
		t.Record(time.Now(), time.Now())

		spy := newSpyLogClient()
		t.Emit(spy)
		Expect(spy.Timers()).To(HaveLen(2))

		spy.ResetTimers()
		t.Emit(spy)
		Expect(spy.Timers()).To(HaveLen(0))
	})

	It("can record timers from multiple threads", func() {
		wg := sync.WaitGroup{}

		t := pulseemitter.NewTimerMetric(
			"some-timer",
			"my-source-id",
			pulseemitter.WithVersion(1, 2),
		)

		numThreads := 16
		numRecordsPerThread := 4
		for i := 0; i < numThreads; i++ {
			wg.Add(1)

			go func() {
				for i := 0; i < numRecordsPerThread; i++ {
					t.Record(time.Now(), time.Now())
					time.Sleep(10 * time.Millisecond)
				}
				wg.Done()
			}()
		}
		wg.Wait()

		spy := newSpyLogClient()
		t.Emit(spy)
		Expect(spy.Timers()).To(HaveLen(numThreads * numRecordsPerThread))
	})
})
