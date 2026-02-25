package rfc5424_test

import (
	"fmt"
	"log"
	"strings"
	"time"

	"code.cloudfoundry.org/go-loggregator/v10/rfc5424"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("RFC5424 Marshaling", func() {

	It("Can marshal and unmarshal", func() {
		var testCases = []struct {
			in       rfc5424.Message
			expected string
		}{
			// RFC-5424 Example 1
			{rfc5424.Message{
				Priority:       34,
				Timestamp:      T("2003-08-24T05:14:15.000003-07:00"),
				Hostname:       "mymachine.example.com",
				AppName:        "su",
				MessageID:      "ID47",
				StructuredData: []rfc5424.StructuredData{},
				Message:        []byte("'su root' failed for lonvick on /dev/pts/8"),
			}, `<34>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com su - ID47 - 'su root' failed for lonvick on /dev/pts/8`},

			// RFC-5424 Example 2
			{rfc5424.Message{
				Priority:       165,
				Timestamp:      T("2003-08-24T05:14:15.000003-07:00"),
				Hostname:       "192.0.2.1",
				AppName:        "myproc",
				ProcessID:      "8710",
				StructuredData: []rfc5424.StructuredData{},
				Message:        []byte("%% It's time to make the do-nuts."),
			}, `<165>1 2003-08-24T05:14:15.000003-07:00 192.0.2.1 myproc 8710 - - %% It's time to make the do-nuts.`},

			// RFC-5424 Example 3
			{rfc5424.Message{
				Priority:  165,
				Timestamp: T("2003-08-24T05:14:15.000003-07:00"),
				Hostname:  "mymachine.example.com",
				AppName:   "evntslog",
				MessageID: "ID47",
				StructuredData: []rfc5424.StructuredData{
					{
						ID: "exampleSDID@32473",
						Parameters: []rfc5424.SDParam{
							{
								Name:  "iut",
								Value: "3",
							},
							{
								Name:  "eventSource",
								Value: "Application",
							},
							{
								Name:  "eventID",
								Value: "1011",
							},
						},
					},
				},
				Message: []byte("An application event log entry..."),
			}, `<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"] An application event log entry...`},

			// RFC-5424 Example 4
			{rfc5424.Message{
				Priority:  165,
				Timestamp: T("2003-08-24T05:14:15.000003-07:00"),
				Hostname:  "mymachine.example.com",
				AppName:   "evntslog",
				MessageID: "ID47",
				StructuredData: []rfc5424.StructuredData{
					{
						ID: "exampleSDID@32473",
						Parameters: []rfc5424.SDParam{
							{
								Name:  "iut",
								Value: "3",
							},
							{
								Name:  "eventSource",
								Value: "Application",
							},
							{
								Name:  "eventID",
								Value: "1011",
							},
						},
					},
					{
						ID: "examplePriority@32473",
						Parameters: []rfc5424.SDParam{
							{
								Name:  "class",
								Value: "high",
							},
						},
					},
				},
			}, `<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47 [exampleSDID@32473 iut="3" eventSource="Application" eventID="1011"][examplePriority@32473 class="high"]`},

			{rfc5424.Message{
				Timestamp: T("2003-08-24T05:14:15.000003-07:00"),
				StructuredData: []rfc5424.StructuredData{
					{
						ID: "x@1",
						Parameters: []rfc5424.SDParam{
							{
								Name:  "class",
								Value: `backslash=\ quote=" right bracket=] left bracket=[`,
							},
						},
					},
				},
			}, `<0>1 2003-08-24T05:14:15.000003-07:00 - - - - [x@1 class="backslash=\\ quote=\" right bracket=\] left bracket=["]`},

			{rfc5424.Message{
				Timestamp:      T("2003-08-24T05:14:15.000003-07:00"),
				StructuredData: []rfc5424.StructuredData{},
			}, `<0>1 2003-08-24T05:14:15.000003-07:00 - - - - -`},

			// UTC TIMESTAMP
			{rfc5424.Message{
				Timestamp:      UTC("2003-08-24T05:14:15.000003Z"),
				StructuredData: []rfc5424.StructuredData{},
				UseUTC:         true,
			}, `<0>1 2003-08-24T05:14:15.000003Z - - - - -`},

			{rfc5424.Message{
				Timestamp: T("2003-08-24T05:14:15.000003-07:00"),
				StructuredData: []rfc5424.StructuredData{
					{
						ID: "x@1",
						Parameters: []rfc5424.SDParam{
							{
								Name:  "",
								Value: "value",
							},
						},
					},
				},
			}, `<0>1 2003-08-24T05:14:15.000003-07:00 - - - - [x@1 ="value"]`},
		}

		for _, tt := range testCases {
			actual, err := tt.in.MarshalBinary()
			Expect(err).To(BeNil())
			Expect(string(actual)).To(Equal(tt.expected))

			m := rfc5424.Message{}
			err = m.UnmarshalBinary(actual)
			if err != nil {
				log.Printf(": %s", string(actual))
				log.Printf(": %#v", m)
			}
			Expect(err).To(BeNil())
			Expect(m).To(Equal(tt.in))
		}
	})

	It("truncates RFC5424 messages", func() {
		longMessage := rfc5424.Message{
			Timestamp:      T("2003-08-24T05:14:15.000003-07:00"),
			Hostname:       strings.Repeat("A", 300),
			AppName:        strings.Repeat("A", 300),
			MessageID:      strings.Repeat("A", 300),
			ProcessID:      strings.Repeat("A", 300),
			StructuredData: []rfc5424.StructuredData{},
		}
		actual, err := longMessage.MarshalBinary()
		Expect(err).To(BeNil())

		expected := fmt.Sprintf(`<0>1 2003-08-24T05:14:15.000003-07:00 %s %s %s %s -`, strings.Repeat("A", 255), strings.Repeat("A", 48), strings.Repeat("A", 128), strings.Repeat("A", 32))
		Expect(string(actual)).To(Equal(expected))
	})

	It("unmarshals valid strings", func() {
		var validStrings = [][]byte{
			[]byte(`<34>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com su X ID47 - msg`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47 [id name="value"]`),
			[]byte(`<165>1 2003-08-24T05:14:15.003-07:00 mymachine.example.com evntslog - ID47 [id name="value"]`),
			[]byte(`<165>1 2003-08-24T05:14:15-07:00 mymachine.example.com evntslog - ID47 [id name="value"]`),
			[]byte(`<165>1 2003-08-24T05:14:15Z mymachine.example.com evntslog - ID47 [id name="value"]`),
			[]byte(`<165>1 2003-08-24T05:14:15.00Z mymachine.example.com evntslog - ID47 [id name="value"]`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003Z mymachine.example.com evntslog - ID47 [id name="value"]`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003+07:00 mymachine.example.com evntslog - ID47 [id name="value"]`),
		}

		for _, actual := range validStrings {
			m := rfc5424.Message{}
			err := m.UnmarshalBinary(actual)
			Expect(err).To(BeNil())
		}
	})

	It("fails to unmarshal invalid strings", func() {
		var invalidStrings = [][]byte{
			[]byte(``),
			[]byte(`<`),
			[]byte(`<3`),
			[]byte(`<34>`),
			[]byte(`<34>1`),
			[]byte(`<34>1 `),
			[]byte(`<34>1 2003-08-24T05:14:15.000003-07:00`),
			[]byte(`<34>1 2003-08-24T05:14:15.000003-07:00`),
			[]byte(`<34>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com`),
			[]byte(`<34>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com su`),
			[]byte(`<34>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com su X`),
			[]byte(`<34>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com su X ID47`),
			[]byte(`<F>1 2003-08-24T05:14:15.000003-07:00 mymachi mymachine.example.com su - ID47 - msg`),
			[]byte(`<34>X 2003-08-24T05:14:15.000003-07:00 mymachine.example.com su - ID47 - msg`),
			[]byte(`<34>1 notATimestamp mymachine.example.com su - ID47 - 'su root' failed for lonvick on /dev/pts/8`),
			[]byte(`>34<1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com su X ID47 - msg`),
			[]byte(`<3499999999999999999999999999999999>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com su X ID47 - msg`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47 `),
			[]byte(`<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47 ]`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47 [`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47 [id`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47 [id name=`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47 [id name="]`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47 [id name="value`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47 [id name="value"`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47 [id name="value"x]`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003-07:00 mymachine.example.com evntslog - ID47 [id name="value\`),
			[]byte(`<165>1 2003-08-24T05:14:15.000003Z+07:00 mymachine.example.com evntslog - ID47 [id name="value"]`),
		}

		for _, actual := range invalidStrings {
			m := rfc5424.Message{}
			err := m.UnmarshalBinary(actual)
			if err == nil {
				log.Printf(": %s", actual)
				log.Printf(": %#v", m)
			}
			Expect(err).NotTo(BeNil())
			Expect(fmt.Sprintf("%s", err)).NotTo(Equal(""))
		}
	})

	It("cannot marshal invalid messages", func() {
		var invalidMessages = []rfc5424.Message{
			{Hostname: "\x7f"},
			{Hostname: "\x20"},
			{AppName: "\x7f"},
			{ProcessID: "\x7f"},
			{MessageID: "\x7f"},
			{
				StructuredData: []rfc5424.StructuredData{
					{
						ID:         "\x20",
						Parameters: []rfc5424.SDParam{{Name: "", Value: "value"}},
					},
				},
			},
			{
				StructuredData: []rfc5424.StructuredData{
					{
						ID:         "\x7f",
						Parameters: []rfc5424.SDParam{{Name: "", Value: "value"}},
					},
				},
			},
			{
				StructuredData: []rfc5424.StructuredData{
					{
						ID:         "foo=bar",
						Parameters: []rfc5424.SDParam{{Name: "", Value: "value"}},
					},
				},
			},
			{
				StructuredData: []rfc5424.StructuredData{
					{
						ID:         "foo[bar]",
						Parameters: []rfc5424.SDParam{{Name: "", Value: "value"}},
					},
				},
			},
			{
				StructuredData: []rfc5424.StructuredData{
					{
						ID:         `foo"bar`,
						Parameters: []rfc5424.SDParam{{Name: "", Value: "value"}},
					},
				},
			},
			{
				StructuredData: []rfc5424.StructuredData{
					{
						ID:         `x@1`,
						Parameters: []rfc5424.SDParam{{Name: "\x7f", Value: "value"}},
					},
				},
			},
			{
				StructuredData: []rfc5424.StructuredData{
					{
						ID:         `x@1`,
						Parameters: []rfc5424.SDParam{{Name: "x", Value: "\xc3\x28"}},
					},
				},
			},
		}

		for i, m := range invalidMessages {
			bin, err := m.MarshalBinary()
			if err == nil {
				log.Printf(": %d", i)
				log.Printf(": %s", string(bin))
				log.Printf(": %#v", m)
			}
			Expect(err).NotTo(BeNil())
			Expect(fmt.Sprintf("%s", err)).NotTo(Equal(""))
		}
	})

	// This test is successful if const allowLongSdNames = true
	It("tests if long attributes are marshaled", func() {
		var message = rfc5424.Message{
			Timestamp: T("2003-08-24T05:14:15.000003-07:00"),
			StructuredData: []rfc5424.StructuredData{
				{
					ID:         "AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA",
					Parameters: []rfc5424.SDParam{{Name: "", Value: "value"}},
				},
			},
		}

		bin, err := message.MarshalBinary()
		Expect(err).To(BeNil())
		Expect(string(bin)).To(Equal("<0>1 2003-08-24T05:14:15.000003-07:00 - - - - [AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA =\"value\"]"))
	})
})

func T(s string) time.Time {
	rv, err := time.Parse(rfc5424.RFC5424TimeOffsetNum, s)
	if err != nil {
		panic(err)
	}
	return rv
}

func UTC(s string) time.Time {
	rv, err := time.Parse(rfc5424.RFC5424TimeOffsetUTC, s)
	if err != nil {
		panic(err)
	}
	return rv
}
