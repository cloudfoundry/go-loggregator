package v2_test

import (
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc/grpclog"

	"testing"
)

func TestV2(t *testing.T) {
	grpclog.SetLogger(log.New(GinkgoWriter, "", 0))
	RegisterFailHandler(Fail)
	RunSpecs(t, "V2 Suite")
}
