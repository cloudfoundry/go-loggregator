package loggregator_test

import (
	"code.cloudfoundry.org/go-loggregator/v10"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("TLS", func() {
	Describe("NewIngressTLSConfig", func() {
		It("returns a valid TLS config with server name 'metron'", func() {
			tlsConf, err := loggregator.NewIngressTLSConfig(
				certs.CA(),
				certs.Cert("test"),
				certs.Key("test"),
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(tlsConf).NotTo(BeNil())
			Expect(tlsConf.ServerName).To(Equal("metron"))
		})

		It("returns an error with an invalid cert path", func() {
			_, err := loggregator.NewIngressTLSConfig(
				certs.CA(),
				"/invalid/cert/path",
				certs.Key("test"),
			)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("NewEgressTLSConfig", func() {
		It("returns a valid TLS config with server name 'reverselogproxy'", func() {
			tlsConf, err := loggregator.NewEgressTLSConfig(
				certs.CA(),
				certs.Cert("test"),
				certs.Key("test"),
			)
			Expect(err).NotTo(HaveOccurred())
			Expect(tlsConf).NotTo(BeNil())
			Expect(tlsConf.ServerName).To(Equal("reverselogproxy"))
		})
	})
})
