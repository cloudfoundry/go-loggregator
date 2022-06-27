package conversion_test

import (
	"code.cloudfoundry.org/go-loggregator/v8/conversion"
	"code.cloudfoundry.org/go-loggregator/v8/rpc/loggregator_v2"
	"github.com/cloudfoundry/sonde-go/events"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gstruct"
	"google.golang.org/protobuf/proto"
)

var _ = Describe("ContainerMetric", func() {
	Context("given a v2 envelope", func() {
		It("converts to a v1 envelope", func() {
			envelope := &loggregator_v2.Envelope{
				SourceId:   "some-id",
				InstanceId: "123",
				Message: &loggregator_v2.Envelope_Gauge{
					Gauge: &loggregator_v2.Gauge{
						Metrics: map[string]*loggregator_v2.GaugeValue{
							"cpu": {
								Unit:  "percentage",
								Value: 11,
							},
							"memory": {
								Unit:  "bytes",
								Value: 13,
							},
							"disk": {
								Unit:  "bytes",
								Value: 15,
							},
							"memory_quota": {
								Unit:  "bytes",
								Value: 17,
							},
							"disk_quota": {
								Unit:  "bytes",
								Value: 19,
							},
						},
					},
				},
			}

			envelopes := conversion.ToV1(envelope)
			Expect(len(envelopes)).To(Equal(1))
			Expect(envelopes[0].GetEventType()).To(Equal(events.Envelope_ContainerMetric))
			Expect(proto.Equal(envelopes[0].GetContainerMetric(), &events.ContainerMetric{
				ApplicationId:    proto.String("some-id"),
				InstanceIndex:    proto.Int32(123),
				CpuPercentage:    proto.Float64(11),
				MemoryBytes:      proto.Uint64(13),
				DiskBytes:        proto.Uint64(15),
				MemoryBytesQuota: proto.Uint64(17),
				DiskBytesQuota:   proto.Uint64(19),
			})).To(BeTrue())
		})

		It("sets InstanceIndex from GaugeValue if present", func() {
			envelope := &loggregator_v2.Envelope{
				InstanceId: "123",
				Message: &loggregator_v2.Envelope_Gauge{
					Gauge: &loggregator_v2.Gauge{
						Metrics: map[string]*loggregator_v2.GaugeValue{
							"instance_index": {
								Unit:  "",
								Value: 19,
							},
							"cpu":          {Unit: "percent"},
							"memory":       {Unit: "percent"},
							"disk":         {Unit: "percent"},
							"memory_quota": {Unit: "bytes"},
							"disk_quota":   {Unit: "bytes"},
						},
					},
				},
			}

			envelopes := conversion.ToV1(envelope)
			Expect(len(envelopes)).To(Equal(1))
			Expect(envelopes[0].GetEventType()).To(Equal(events.Envelope_ContainerMetric))
			Expect(proto.Equal(envelopes[0].GetContainerMetric(), &events.ContainerMetric{
				ApplicationId:    proto.String(""),
				InstanceIndex:    proto.Int32(19),
				CpuPercentage:    proto.Float64(0),
				MemoryBytes:      proto.Uint64(0),
				DiskBytes:        proto.Uint64(0),
				MemoryBytesQuota: proto.Uint64(0),
				DiskBytesQuota:   proto.Uint64(0),
			})).To(BeTrue())
		})

		DescribeTable("it is resilient to malformed envelopes", func(v2e *loggregator_v2.Envelope) {
			Expect(conversion.ToV1(v2e)).To(HaveLen(0))
		},
			Entry("bare envelope", &loggregator_v2.Envelope{}),
			Entry("with empty fields", &loggregator_v2.Envelope{
				Message: &loggregator_v2.Envelope_Gauge{
					Gauge: &loggregator_v2.Gauge{
						Metrics: map[string]*loggregator_v2.GaugeValue{
							"cpu":          nil,
							"memory":       nil,
							"disk":         nil,
							"memory_quota": nil,
							"disk_quota":   nil,
						},
					},
				},
			}),
		)
	})

	Context("given a v1 envelope", func() {
		var (
			v1Envelope      *events.Envelope
			expectedMessage *loggregator_v2.Envelope_Gauge
		)

		BeforeEach(func() {
			v1Envelope = &events.Envelope{
				Origin:     proto.String("an-origin"),
				Deployment: proto.String("a-deployment"),
				Job:        proto.String("a-job"),
				Index:      proto.String("an-index"),
				Ip:         proto.String("an-ip"),
				Timestamp:  proto.Int64(1234),
				EventType:  events.Envelope_ContainerMetric.Enum(),
				ContainerMetric: &events.ContainerMetric{
					ApplicationId:    proto.String("some-id"),
					InstanceIndex:    proto.Int32(123),
					CpuPercentage:    proto.Float64(11),
					MemoryBytes:      proto.Uint64(13),
					DiskBytes:        proto.Uint64(15),
					MemoryBytesQuota: proto.Uint64(17),
					DiskBytesQuota:   proto.Uint64(19),
				},
				Tags: map[string]string{
					"custom_tag": "custom-value",
				},
			}
			expectedMessage = &loggregator_v2.Envelope_Gauge{
				Gauge: &loggregator_v2.Gauge{
					Metrics: map[string]*loggregator_v2.GaugeValue{
						"cpu": {
							Unit:  "percentage",
							Value: 11,
						},
						"memory": {
							Unit:  "bytes",
							Value: 13,
						},
						"disk": {
							Unit:  "bytes",
							Value: 15,
						},
						"memory_quota": {
							Unit:  "bytes",
							Value: 17,
						},
						"disk_quota": {
							Unit:  "bytes",
							Value: 19,
						},
					},
				},
			}
		})

		Context("using deprecated tags", func() {
			It("converts to a v2 envelope using DeprecatedTags", func() {
				Expect(*conversion.ToV2(v1Envelope, false)).To(MatchFields(IgnoreExtras, Fields{
					"Timestamp":  Equal(int64(1234)),
					"SourceId":   Equal("some-id"),
					"Message":    Equal(expectedMessage),
					"InstanceId": Equal("123"),
					"DeprecatedTags": Equal(map[string]*loggregator_v2.Value{
						"origin":     {Data: &loggregator_v2.Value_Text{Text: "an-origin"}},
						"deployment": {Data: &loggregator_v2.Value_Text{Text: "a-deployment"}},
						"job":        {Data: &loggregator_v2.Value_Text{Text: "a-job"}},
						"index":      {Data: &loggregator_v2.Value_Text{Text: "an-index"}},
						"ip":         {Data: &loggregator_v2.Value_Text{Text: "an-ip"}},
						"__v1_type":  {Data: &loggregator_v2.Value_Text{Text: "ContainerMetric"}},
						"custom_tag": {Data: &loggregator_v2.Value_Text{Text: "custom-value"}},
					}),
					"Tags": BeNil(),
				}))
			})

			It("sets the source ID to deployment/job when App ID is missing", func() {
				localV1Envelope := &events.Envelope{
					Deployment:      proto.String("some-deployment"),
					Job:             proto.String("some-job"),
					EventType:       events.Envelope_ContainerMetric.Enum(),
					ContainerMetric: &events.ContainerMetric{},
				}

				expectedV2Envelope := &loggregator_v2.Envelope{
					SourceId: "some-deployment/some-job",
				}

				converted := conversion.ToV2(localV1Envelope, false)
				Expect(converted.GetSourceId()).To(Equal(expectedV2Envelope.GetSourceId()))
			})
		})

		Context("using preferred tags", func() {
			It("converts to a v2 envelope using Tags", func() {
				Expect(*conversion.ToV2(v1Envelope, true)).To(MatchFields(IgnoreExtras, Fields{
					"Timestamp":      Equal(int64(1234)),
					"SourceId":       Equal("some-id"),
					"Message":        Equal(expectedMessage),
					"InstanceId":     Equal("123"),
					"DeprecatedTags": BeNil(),
					"Tags": Equal(map[string]string{
						"origin":     "an-origin",
						"deployment": "a-deployment",
						"job":        "a-job",
						"index":      "an-index",
						"ip":         "an-ip",
						"__v1_type":  "ContainerMetric",
						"custom_tag": "custom-value",
					}),
				}))
			})
		})
	})
})
