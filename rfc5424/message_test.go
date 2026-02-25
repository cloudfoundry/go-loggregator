package rfc5424_test

import (
	"code.cloudfoundry.org/go-loggregator/v10/rfc5424"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RFC5424 Message", func() {
	It("adds datum to a message", func() {
		m := rfc5424.Message{}
		m.AddDatum("id", "name", "value")
		Expect(m).To(Equal(rfc5424.Message{
			StructuredData: []rfc5424.StructuredData{
				{
					ID: "id",
					Parameters: []rfc5424.SDParam{
						{"name", "value"},
					},
				},
			},
		}))

		m.AddDatum("id2", "name", "value")
		Expect(m).To(Equal(rfc5424.Message{
			StructuredData: []rfc5424.StructuredData{
				{
					ID: "id",
					Parameters: []rfc5424.SDParam{
						{"name", "value"},
					},
				},
				{
					ID: "id2",
					Parameters: []rfc5424.SDParam{
						{"name", "value"},
					},
				},
			},
		}))

		m.AddDatum("id", "name2", "value2")
		Expect(m).To(Equal(rfc5424.Message{
			StructuredData: []rfc5424.StructuredData{
				{
					ID: "id",
					Parameters: []rfc5424.SDParam{
						{"name", "value"},
						{"name2", "value2"},
					},
				},
				{
					ID: "id2",
					Parameters: []rfc5424.SDParam{
						{"name", "value"},
					},
				},
			},
		}))

		m.AddDatum("id", "name", "value3")
		Expect(m).To(Equal(rfc5424.Message{
			StructuredData: []rfc5424.StructuredData{
				{
					ID: "id",
					Parameters: []rfc5424.SDParam{
						{"name", "value"},
						{"name2", "value2"},
						{"name", "value3"},
					},
				},
				{
					ID: "id2",
					Parameters: []rfc5424.SDParam{
						{"name", "value"},
					},
				},
			},
		}))
	})
})
