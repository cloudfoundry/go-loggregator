package v2_test

import (
	"code.cloudfoundry.org/go-loggregator/v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("NewTLSConfig", func() {
	It("works with valid certs", func() {
		_, err := v2.NewTLSConfig("../fixtures/CA.crt", "../fixtures/client.crt", "../fixtures/client.key")

		Expect(err).NotTo(HaveOccurred())
	})

	It("errors with invalid cert", func() {
		_, err := v2.NewTLSConfig("../fixtures/CA.crt", "invalid", "../fixtures/client.key")
		Expect(err).To(HaveOccurred(), "client didn't return an error")
	})

	It("errors with invalid key", func() {
		_, err := v2.NewTLSConfig("../fixtures/CA.crt", "../fixtures/client.crt", "invalid")
		Expect(err).To(HaveOccurred(), "client didn't return an error")
	})

	It("errors with invalid CA cert", func() {
		_, err := v2.NewTLSConfig("invalid", "../fixtures/client.crt", "../fixtures/client.key")
		Expect(err).To(HaveOccurred(), "client didn't return an error")
	})
})
