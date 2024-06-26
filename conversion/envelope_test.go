package conversion_test

import (
	"fmt"

	"code.cloudfoundry.org/go-loggregator/v10/conversion"
	"code.cloudfoundry.org/go-loggregator/v10/rpc/loggregator_v2"
	"github.com/cloudfoundry/sonde-go/events"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/protobuf/proto"
)

var _ = Describe("Envelope", func() {
	Context("given a v2 envelope", func() {
		It("sets v1 specific properties", func() {
			envelope := &loggregator_v2.Envelope{
				Timestamp: 99,
				DeprecatedTags: map[string]*loggregator_v2.Value{
					"origin":         {Data: &loggregator_v2.Value_Text{Text: "origin"}},
					"deployment":     {Data: &loggregator_v2.Value_Text{Text: "deployment"}},
					"job":            {Data: &loggregator_v2.Value_Text{Text: "job"}},
					"index":          {Data: &loggregator_v2.Value_Text{Text: "index"}},
					"ip":             {Data: &loggregator_v2.Value_Text{Text: "ip"}},
					"random_text":    {Data: &loggregator_v2.Value_Text{Text: "random_text"}},
					"random_int":     {Data: &loggregator_v2.Value_Integer{Integer: 123}},
					"random_decimal": {Data: &loggregator_v2.Value_Decimal{Decimal: 123}},
				},
				Message: &loggregator_v2.Envelope_Log{Log: &loggregator_v2.Log{}},
			}

			envelopes := conversion.ToV1(envelope)
			Expect(len(envelopes)).To(Equal(1))
			oldEnvelope := envelopes[0]
			Expect(oldEnvelope.GetOrigin()).To(Equal("origin"))
			Expect(oldEnvelope.GetEventType()).To(Equal(events.Envelope_LogMessage))
			Expect(oldEnvelope.GetTimestamp()).To(Equal(int64(99)))
			Expect(oldEnvelope.GetDeployment()).To(Equal("deployment"))
			Expect(oldEnvelope.GetJob()).To(Equal("job"))
			Expect(oldEnvelope.GetIndex()).To(Equal("index"))
			Expect(oldEnvelope.GetIp()).To(Equal("ip"))
			Expect(oldEnvelope.Tags).To(HaveKeyWithValue("random_text", "random_text"))
			Expect(oldEnvelope.Tags).To(HaveKeyWithValue("random_int", "123"))
			Expect(oldEnvelope.Tags).To(HaveKeyWithValue("random_decimal", fmt.Sprintf("%f", 123.0)))
		})

		It("rejects empty tags", func() {
			envelope := &loggregator_v2.Envelope{
				DeprecatedTags: map[string]*loggregator_v2.Value{
					"foo": {Data: &loggregator_v2.Value_Text{Text: "bar"}},
					"baz": nil,
				},
				Message: &loggregator_v2.Envelope_Log{Log: &loggregator_v2.Log{}},
			}

			envelopes := conversion.ToV1(envelope)
			Expect(len(envelopes)).To(Equal(1))
			oldEnvelope := envelopes[0]
			Expect(oldEnvelope.Tags).To(Equal(map[string]string{
				"foo": "bar",
			}))
		})

		It("reads non-text v2 tags", func() {
			envelope := &loggregator_v2.Envelope{
				DeprecatedTags: map[string]*loggregator_v2.Value{
					"foo": {Data: &loggregator_v2.Value_Integer{Integer: 99}},
				},
				Message: &loggregator_v2.Envelope_Log{Log: &loggregator_v2.Log{}},
			}

			envelopes := conversion.ToV1(envelope)
			Expect(len(envelopes)).To(Equal(1))
			Expect(envelopes[0].GetTags()).To(HaveKeyWithValue("foo", "99"))
		})

		It("uses non-deprecated v2 tags", func() {
			envelope := &loggregator_v2.Envelope{
				Tags: map[string]string{
					"foo": "bar",
				},
				Message: &loggregator_v2.Envelope_Log{Log: &loggregator_v2.Log{}},
			}

			envelopes := conversion.ToV1(envelope)
			Expect(len(envelopes)).To(Equal(1))
			Expect(envelopes[0].GetTags()).To(HaveKeyWithValue("foo", "bar"))
		})
	})

	Context("given a v1 envelope", func() {
		It("sets v2 specific properties", func() {
			v1Envelope := &events.Envelope{
				Timestamp:  proto.Int64(99),
				Origin:     proto.String("origin-value"),
				Deployment: proto.String("some-deployment"),
				Job:        proto.String("some-job"),
				Index:      proto.String("some-index"),
				Ip:         proto.String("some-ip"),
				Tags: map[string]string{
					"random-tag": "random-value",
				},
			}

			expectedV2Envelope := &loggregator_v2.Envelope{
				Timestamp: 99,
				SourceId:  "some-deployment/some-job",
				DeprecatedTags: map[string]*loggregator_v2.Value{
					"random-tag": ValueText("random-value"),
					"origin":     ValueText("origin-value"),
					"deployment": ValueText("some-deployment"),
					"job":        ValueText("some-job"),
					"index":      ValueText("some-index"),
					"ip":         ValueText("some-ip"),
				},
			}

			converted := conversion.ToV2(v1Envelope, false)

			Expect(converted.SourceId).To(Equal(expectedV2Envelope.SourceId))
			Expect(converted.Timestamp).To(Equal(expectedV2Envelope.Timestamp))
			Expect(converted.DeprecatedTags["random-tag"]).To(Equal(expectedV2Envelope.DeprecatedTags["random-tag"]))
			Expect(converted.DeprecatedTags["origin"]).To(Equal(expectedV2Envelope.DeprecatedTags["origin"]))
			Expect(converted.DeprecatedTags["deployment"]).To(Equal(expectedV2Envelope.DeprecatedTags["deployment"]))
			Expect(converted.DeprecatedTags["job"]).To(Equal(expectedV2Envelope.DeprecatedTags["job"]))
			Expect(converted.DeprecatedTags["index"]).To(Equal(expectedV2Envelope.DeprecatedTags["index"]))
			Expect(converted.DeprecatedTags["ip"]).To(Equal(expectedV2Envelope.DeprecatedTags["ip"]))
		})

		It("sets non-deprecated tags", func() {
			v1 := &events.Envelope{
				Timestamp:  proto.Int64(99),
				Origin:     proto.String("origin-value"),
				Deployment: proto.String("some-deployment"),
				Job:        proto.String("some-job"),
				Index:      proto.String("some-index"),
				Ip:         proto.String("some-ip"),
				Tags: map[string]string{
					"random-tag": "random-value",
					"origin":     "origin-value",
					"deployment": "some-deployment",
					"job":        "some-job",
					"index":      "some-index",
					"ip":         "some-ip",
				},
			}
			expected := proto.Clone(v1).(*events.Envelope)

			converted := conversion.ToV2(v1, true)

			Expect(converted.Tags["random-tag"]).To(Equal(expected.Tags["random-tag"]))
			Expect(converted.Tags["origin"]).To(Equal(expected.Tags["origin"]))
			Expect(converted.Tags["deployment"]).To(Equal(expected.Tags["deployment"]))
			Expect(converted.Tags["job"]).To(Equal(expected.Tags["job"]))
			Expect(converted.Tags["index"]).To(Equal(expected.Tags["index"]))
			Expect(converted.Tags["ip"]).To(Equal(expected.Tags["ip"]))
		})
	})
})
