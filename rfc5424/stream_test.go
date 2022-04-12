package rfc5424

import (
	"bytes"

	. "gopkg.in/check.v1"
)

var _ = Suite(&StreamTest{})

type StreamTest struct {
}

func (s *StreamTest) TestCanReadAndWrite(c *C) {
	stream := bytes.Buffer{}
	for i := 0; i < 4; i++ {
		m := Message{Priority: Priority(i), Timestamp: T("2003-08-24T05:14:15.000003-07:00")}
		nbytes, err := m.WriteTo(&stream)
		c.Assert(err, IsNil)
		c.Assert(nbytes, Equals, int64(50))
	}

	c.Assert(stream.String(), Equals,
		`47 <0>1 2003-08-24T05:14:15.000003-07:00 - - - - -`+
			`47 <1>1 2003-08-24T05:14:15.000003-07:00 - - - - -`+
			`47 <2>1 2003-08-24T05:14:15.000003-07:00 - - - - -`+
			`47 <3>1 2003-08-24T05:14:15.000003-07:00 - - - - -`)

	for i := 0; i < 4; i++ {
		m := Message{Priority: Priority(i << 3)}
		nbytes, err := m.ReadFrom(&stream)
		c.Assert(err, IsNil)
		c.Assert(nbytes, Equals, int64(50))
		c.Assert(m, DeepEquals, Message{Priority: Priority(i),
			Timestamp:      T("2003-08-24T05:14:15.000003-07:00"),
			StructuredData: []StructuredData{}})
	}
}

func (s *StreamTest) TestUtcTimestamps(c *C) {
	stream := bytes.Buffer{}
	for i := 0; i < 4; i++ {
		m := Message{Priority: Priority(i), Timestamp: T("2003-08-24T05:14:15.000003+00:00"), UseUTC: true}
		nbytes, err := m.WriteTo(&stream)
		c.Assert(err, IsNil)
		c.Assert(nbytes, Equals, int64(45))
	}

	c.Assert(stream.String(), Equals,
		`42 <0>1 2003-08-24T05:14:15.000003Z - - - - -`+
			`42 <1>1 2003-08-24T05:14:15.000003Z - - - - -`+
			`42 <2>1 2003-08-24T05:14:15.000003Z - - - - -`+
			`42 <3>1 2003-08-24T05:14:15.000003Z - - - - -`)

	for i := 0; i < 4; i++ {
		m := Message{Priority: Priority(i << 3)}
		nbytes, err := m.ReadFrom(&stream)
		c.Assert(err, IsNil)
		c.Assert(nbytes, Equals, int64(45))
		c.Assert(m, DeepEquals, Message{Priority: Priority(i),
			Timestamp:      UTC("2003-08-24T05:14:15.000003Z"),
			StructuredData: []StructuredData{},
			UseUTC:         true,
		})
	}
}

func (s *StreamTest) TestRejectsInvalidStream(c *C) {
	stream := bytes.NewBufferString(`99 <0>1 2003-08-24T05:14:15.000003-07:00 - - - - -`)
	for i := 0; i < 4; i++ {
		m := Message{Priority: Priority(i << 3)}
		_, err := m.ReadFrom(stream)
		c.Assert(err, Not(IsNil))
	}
}

func (s *StreamTest) TestRejectsInvalidStream2(c *C) {
	stream := bytes.NewBufferString(`0 <0>1 2003-08-24T05:14:15.000003-07:00 - - - - -`)
	for i := 0; i < 4; i++ {
		m := Message{Priority: Priority(i << 3)}
		_, err := m.ReadFrom(stream)
		c.Assert(err, Not(IsNil))
	}
}
