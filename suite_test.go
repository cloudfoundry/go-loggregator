package loggregator_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestLoggregatorV2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "go-loggregator compilation check")
}
