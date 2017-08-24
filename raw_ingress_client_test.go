package loggregator_test

import (
	"code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RawIngressClient", func() {
	var (
		client    *loggregator.RawIngressClient
		receivers chan loggregator_v2.Ingress_BatchSenderServer
		server    *testIngressServer
	)

	BeforeEach(func() {
		var err error
		server, err = newTestIngressServer(
			fixture("server.crt"),
			fixture("server.key"),
			fixture("CA.crt"),
		)
		Expect(err).NotTo(HaveOccurred())

		err = server.start()
		Expect(err).NotTo(HaveOccurred())
		receivers = server.receivers()

		tlsConfig, err := loggregator.NewIngressTLSConfig(
			fixture("CA.crt"),
			fixture("client.crt"),
			fixture("client.key"),
		)
		Expect(err).NotTo(HaveOccurred())

		client, err = loggregator.NewRawIngressClient(
			server.addr(),
			tlsConfig,
		)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		server.stop()
	})

	It("reconnects when the server goes away and comes back", func() {
		client.Emit([]*loggregator_v2.Envelope{
			{
				Timestamp: 1,
			},
			{
				Timestamp: 2,
			},
		})

		var recv loggregator_v2.Ingress_BatchSenderServer
		Eventually(receivers, 10).Should(Receive(&recv))
		envBatch, err := recv.Recv()

		Expect(err).NotTo(HaveOccurred())
		Expect(envBatch.Batch).To(HaveLen(2))

		server.stop()
		Eventually(server.start).Should(Succeed())

		Consistently(receivers).Should(BeEmpty())

		go func() {
			for {
				client.Emit([]*loggregator_v2.Envelope{
					{
						Timestamp: 3,
					},
				})
			}
		}()

		Eventually(receivers, 10).Should(Receive(&recv))
		envBatch, err = recv.Recv()

		Expect(err).NotTo(HaveOccurred())
		Expect(envBatch.Batch).ToNot(BeEmpty())
	})
})
