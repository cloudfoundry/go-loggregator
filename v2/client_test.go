package v2_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/go-loggregator/internal/loggregator_v2"
	"code.cloudfoundry.org/go-loggregator/runtimeemitter"
	"code.cloudfoundry.org/go-loggregator/v2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GrpcClient", func() {
	var (
		client    *v2.Client
		clientErr error
		receivers chan loggregator_v2.Ingress_BatchSenderServer
		server    *TestServer
	)

	BeforeEach(func() {
		var err error
		server, err = NewTestServer("fixtures/metron.crt", "fixtures/metron.key", "fixtures/CA.crt")
		Expect(err).NotTo(HaveOccurred())

		err = server.Start()
		Expect(err).NotTo(HaveOccurred())

		receivers = server.Receivers()

		tlsConfig, err := v2.NewTLSConfig(
			"fixtures/CA.crt",
			"fixtures/client.crt",
			"fixtures/client.key",
		)
		Expect(err).NotTo(HaveOccurred())

		client, clientErr = v2.NewClient(
			tlsConfig,
			v2.WithPort(server.Port()),
			v2.WithBatchFlushInterval(50*time.Millisecond),
			v2.WithTag("origin", "some-origin"),
		)
	})

	AfterEach(func() {
		server.Stop()
	})

	It("sends in batches", func() {
		for i := 0; i < 10; i++ {
			client.EmitLog(
				"message",
				v2.WithAppInfo("app-id", "source-type", "source-instance"),
				v2.WithStdout(),
			)
			time.Sleep(10 * time.Millisecond)
		}

		Eventually(func() int {
			b, err := getBatch(receivers)
			if err != nil {
				return 0
			}

			return len(b.Batch)
		}).Should(BeNumerically(">", 1))
	})

	It("sends app logs", func() {
		client.EmitLog(
			"message",
			v2.WithAppInfo("app-id", "source-type", "source-instance"),
			v2.WithStdout(),
		)
		env, err := getEnvelopeAt(receivers, 0)
		Expect(err).NotTo(HaveOccurred())

		Expect(env.Tags["source_instance"].GetText()).To(Equal("source-instance"))
		Expect(env.Tags["origin"].GetText()).To(Equal("some-origin"))
		Expect(env.SourceId).To(Equal("app-id"))
		Expect(env.InstanceId).To(Equal("source-instance"))

		ts := time.Unix(0, env.Timestamp)
		Expect(ts).Should(BeTemporally("~", time.Now(), time.Second))
		log := env.GetLog()
		Expect(log).NotTo(BeNil())
		Expect(log.Payload).To(Equal([]byte("message")))
		Expect(log.Type).To(Equal(loggregator_v2.Log_OUT))
	})

	It("sends app error logs", func() {
		client.EmitLog(
			"message",
			v2.WithAppInfo("app-id", "source-type", "source-instance"),
		)

		env, err := getEnvelopeAt(receivers, 0)
		Expect(err).NotTo(HaveOccurred())

		Expect(env.Tags["source_instance"].GetText()).To(Equal("source-instance"))
		Expect(env.SourceId).To(Equal("app-id"))
		Expect(env.InstanceId).To(Equal("source-instance"))

		ts := time.Unix(0, env.Timestamp)
		Expect(ts).Should(BeTemporally("~", time.Now(), time.Second))
		log := env.GetLog()
		Expect(log).NotTo(BeNil())
		Expect(log.Payload).To(Equal([]byte("message")))
		Expect(log.Type).To(Equal(loggregator_v2.Log_ERR))
	})

	It("sends app metrics", func() {
		client.EmitGauge(
			v2.WithGaugeValue("name-a", 1, "unit-a"),
			v2.WithGaugeValue("name-b", 2, "unit-b"),
			v2.WithGaugeTags(map[string]string{"some-tag": "some-tag-value"}),
			v2.WithGaugeAppInfo("app-id"),
		)

		env, err := getEnvelopeAt(receivers, 0)
		Expect(err).NotTo(HaveOccurred())

		ts := time.Unix(0, env.Timestamp)
		Expect(ts).Should(BeTemporally("~", time.Now(), time.Second))
		metrics := env.GetGauge()
		Expect(metrics).NotTo(BeNil())
		Expect(env.SourceId).To(Equal("app-id"))
		Expect(metrics.GetMetrics()).To(HaveLen(2))
		Expect(metrics.GetMetrics()["name-a"].Value).To(Equal(1.0))
		Expect(metrics.GetMetrics()["name-b"].Value).To(Equal(2.0))
		Expect(env.Tags["some-tag"].GetText()).To(Equal("some-tag-value"))
		Expect(env.Tags["origin"].GetText()).To(Equal("some-origin"))
	})

	It("reconnects when the server goes away and comes back", func() {
		client.EmitLog(
			"message",
			v2.WithAppInfo("app-id", "source-type", "source-instance"),
		)

		envBatch, err := getBatch(receivers)
		Expect(err).NotTo(HaveOccurred())
		Expect(envBatch.Batch).To(HaveLen(1))

		server.Stop()
		Expect(server.Start()).To(Succeed())

		Consistently(receivers).Should(BeEmpty())

		for i := 0; i < 200; i++ {
			client.EmitLog(
				"message",
				v2.WithAppInfo("app-id", "source-type", "source-instance"),
			)
		}

		envBatch, err = getBatch(receivers)
		Expect(err).NotTo(HaveOccurred())
		Expect(envBatch.Batch).ToNot(BeEmpty())
	})

	It("works with the runtime emitter", func() {
		// This test is to ensure that the v2 client satisfies the
		// runtimeemitter.Sender interface. If it does not satisfy the
		// runtimeemitter.Sender interface this test will force a compile time
		// error.
		runtimeemitter.New(client)
	})
})

func getBatch(receivers chan loggregator_v2.Ingress_BatchSenderServer) (*loggregator_v2.EnvelopeBatch, error) {
	var recv loggregator_v2.Ingress_BatchSenderServer
	Eventually(receivers, 10).Should(Receive(&recv))

	return recv.Recv()
}

func getEnvelopeAt(receivers chan loggregator_v2.Ingress_BatchSenderServer, idx int) (*loggregator_v2.Envelope, error) {
	envBatch, err := getBatch(receivers)
	if err != nil {
		return nil, err
	}

	if len(envBatch.Batch) < 1 {
		return nil, errors.New("no envelopes")
	}

	return envBatch.Batch[idx], nil
}
