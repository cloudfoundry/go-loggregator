package conversion_test

import (
	v2 "code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"code.cloudfoundry.org/go-loggregator/v1/conversion"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("Converting Instance IDs", func() {
	Describe("LogMessage", func() {
		It("writes into the v1 source_instance field when converting to v1", func() {
			v2Envelope := &v2.Envelope{
				InstanceId: "test-source-instance",
				Message: &v2.Envelope_Log{
					Log: &v2.Log{
						Payload: []byte("Hello World"),
						Type:    v2.Log_OUT,
					},
				},
			}
			envelopes := conversion.ToV1(v2Envelope)
			Expect(len(envelopes)).To(Equal(1))
			Expect(*envelopes[0].LogMessage.SourceInstance).To(Equal("test-source-instance"))
		})
	})

	Describe("HttpStartStop", func() {
		It("writes into the v1 instance_index field when converting to v1", func() {
			v2Envelope := &v2.Envelope{
				InstanceId: "1234",
				Message: &v2.Envelope_Timer{
					Timer: &v2.Timer{},
				},
			}
			envelopes := conversion.ToV1(v2Envelope)
			Expect(len(envelopes)).To(Equal(1))
			Expect(*envelopes[0].HttpStartStop.InstanceIndex).To(Equal(int32(1234)))
		})

		It("writes 0 into the v1 instance_index field if instance_id is not an int", func() {
			v2Envelope := &v2.Envelope{
				InstanceId: "garbage",
				Message: &v2.Envelope_Timer{
					Timer: &v2.Timer{},
				},
			}
			envelopes := conversion.ToV1(v2Envelope)
			Expect(len(envelopes)).To(Equal(1))
			Expect(*envelopes[0].HttpStartStop.InstanceIndex).To(Equal(int32(0)))
		})
	})

	Describe("ContainerMetric", func() {
		It("writes into the v1 instance_index field when converting to v1", func() {
			v2Envelope := &v2.Envelope{
				InstanceId: "4321",
				Message: &v2.Envelope_Gauge{
					Gauge: &v2.Gauge{
						Metrics: map[string]*v2.GaugeValue{
							"cpu":          {},
							"memory":       {},
							"disk":         {},
							"memory_quota": {},
							"disk_quota":   {},
						},
					},
				},
			}
			envelopes := conversion.ToV1(v2Envelope)
			Expect(len(envelopes)).To(Equal(1))
			Expect(*envelopes[0].ContainerMetric.InstanceIndex).To(Equal(int32(4321)))
		})

		It("writes 0 into the v1 instance_index field if instance_id is not an int", func() {
			v2Envelope := &v2.Envelope{
				InstanceId: "garbage",
				Message: &v2.Envelope_Gauge{
					Gauge: &v2.Gauge{
						Metrics: map[string]*v2.GaugeValue{
							"cpu":          {},
							"memory":       {},
							"disk":         {},
							"memory_quota": {},
							"disk_quota":   {},
						},
					},
				},
			}
			envelopes := conversion.ToV1(v2Envelope)
			Expect(len(envelopes)).To(Equal(1))
			Expect(*envelopes[0].ContainerMetric.InstanceIndex).To(Equal(int32(0)))
		})
	})

	Describe("CounterEvent and ValueMetric", func() {
		DescribeTable("writes into the v1 instance_id tag when converting to v1", func(v2Envelope *v2.Envelope) {
			envelopes := conversion.ToV1(v2Envelope)
			Expect(len(envelopes)).To(Equal(1))
			Expect(envelopes[0].Tags["instance_id"]).To(Equal("test-source-instance"))
		},
			Entry("CounterEvent", &v2.Envelope{
				InstanceId: "test-source-instance",
				Message: &v2.Envelope_Counter{
					Counter: &v2.Counter{},
				},
			}),
			Entry("ValueMetric", &v2.Envelope{
				InstanceId: "test-source-instance",
				Message: &v2.Envelope_Gauge{
					Gauge: &v2.Gauge{
						Metrics: map[string]*v2.GaugeValue{
							"some-metric": {
								Unit:  "test",
								Value: 123.4,
							},
						},
					},
				},
			}),
		)
	})
})
