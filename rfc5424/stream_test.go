package rfc5424_test

import (
	"bytes"

	"code.cloudfoundry.org/go-loggregator/v10/rfc5424"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RFC5424 Stream", func() {
	It("can read and write messages", func() {
		stream := bytes.Buffer{}
		for i := 0; i < 4; i++ {
			m := rfc5424.Message{Priority: rfc5424.Priority(i), Timestamp: T("2003-08-24T05:14:15.000003-07:00")}
			nbytes, err := m.WriteTo(&stream)
			Expect(err).To(BeNil())
			Expect(nbytes).To(Equal(int64(50)))
		}

		Expect(stream.String()).To(Equal(
			`47 <0>1 2003-08-24T05:14:15.000003-07:00 - - - - -` +
				`47 <1>1 2003-08-24T05:14:15.000003-07:00 - - - - -` +
				`47 <2>1 2003-08-24T05:14:15.000003-07:00 - - - - -` +
				`47 <3>1 2003-08-24T05:14:15.000003-07:00 - - - - -`,
		))

		for i := 0; i < 4; i++ {
			m := rfc5424.Message{Priority: rfc5424.Priority(i << 3)}
			nbytes, err := m.ReadFrom(&stream)
			Expect(err).To(BeNil())
			Expect(nbytes).To(Equal(int64(50)))
			Expect(m).To(Equal(rfc5424.Message{
				Priority:       rfc5424.Priority(i),
				Timestamp:      T("2003-08-24T05:14:15.000003-07:00"),
				StructuredData: []rfc5424.StructuredData{}},
			))
		}
	})

	It("can use UTC Timestamps", func() {
		stream := bytes.Buffer{}
		for i := 0; i < 4; i++ {
			m := rfc5424.Message{Priority: rfc5424.Priority(i), Timestamp: T("2003-08-24T05:14:15.000003+00:00"), UseUTC: true}
			nbytes, err := m.WriteTo(&stream)
			Expect(err).To(BeNil())
			Expect(nbytes).To(Equal(int64(45)))
		}

		Expect(stream.String()).To(Equal(
			`42 <0>1 2003-08-24T05:14:15.000003Z - - - - -` +
				`42 <1>1 2003-08-24T05:14:15.000003Z - - - - -` +
				`42 <2>1 2003-08-24T05:14:15.000003Z - - - - -` +
				`42 <3>1 2003-08-24T05:14:15.000003Z - - - - -`,
		))

		for i := 0; i < 4; i++ {
			m := rfc5424.Message{Priority: rfc5424.Priority(i << 3)}
			nbytes, err := m.ReadFrom(&stream)
			Expect(err).To(BeNil())
			Expect(nbytes).To(Equal(int64(45)))

			Expect(m).To(Equal(rfc5424.Message{
				Priority:       rfc5424.Priority(i),
				Timestamp:      UTC("2003-08-24T05:14:15.000003Z"),
				StructuredData: []rfc5424.StructuredData{},
				UseUTC:         true,
			}))
		}
	})

	It("rejects invalid streams", func() {
		stream := bytes.NewBufferString(`99 <0>1 2003-08-24T05:14:15.000003-07:00 - - - - -`)
		for i := 0; i < 4; i++ {
			m := rfc5424.Message{Priority: rfc5424.Priority(i << 3)}
			_, err := m.ReadFrom(stream)
			Expect(err).NotTo(BeNil())
		}
	})

	It("rejects invalid streams 2", func() {
		stream := bytes.NewBufferString(`0 <0>1 2003-08-24T05:14:15.000003-07:00 - - - - -`)
		for i := 0; i < 4; i++ {
			m := rfc5424.Message{Priority: rfc5424.Priority(i << 3)}
			_, err := m.ReadFrom(stream)
			Expect(err).NotTo(BeNil())
		}
	})
})
