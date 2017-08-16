package conversion_test

import (
	"testing"

	v2 "code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestConversion(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Conversion Suite")
}

func ValueText(s string) *v2.Value {
	return &v2.Value{&v2.Value_Text{Text: s}}
}

func ValueInteger(i int64) *v2.Value {
	return &v2.Value{&v2.Value_Integer{Integer: i}}
}
