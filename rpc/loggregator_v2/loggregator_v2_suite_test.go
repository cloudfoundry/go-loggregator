package loggregator_v2_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestLoggregatorV2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "LoggregatorV2 Suite")
}
