package loggregator_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/go-loggregator"
	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"code.cloudfoundry.org/go-loggregator/runtimeemitter"
	"golang.org/x/net/context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"
)

var _ = Describe("IngressClient", func() {
	var (
		client *loggregator.IngressClient
		server *testIngressServer
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

		client = buildIngressClient(server.addr, 50*time.Millisecond)
	})

	AfterEach(func() {
		server.stop()
	})

	It("sends in batches", func() {
		for i := 0; i < 10; i++ {
			client.EmitLog(
				"message",
				loggregator.WithAppInfo("app-id", "source-type", "source-instance"),
				loggregator.WithStdout(),
			)
			time.Sleep(10 * time.Millisecond)
		}

		Eventually(func() int {
			var recv loggregator_v2.Ingress_BatchSenderServer
			Eventually(server.receivers, 10).Should(Receive(&recv))

			b, err := recv.Recv()
			if err != nil {
				return 0
			}

			return len(b.Batch)
		}).Should(BeNumerically(">", 1))
	})

	It("sends app logs", func() {
		client.EmitLog(
			"message",
			loggregator.WithAppInfo("app-id", "source-type", "source-instance"),
			loggregator.WithStdout(),
		)
		env, err := getEnvelopeAt(server.receivers, 0)
		Expect(err).NotTo(HaveOccurred())

		Expect(env.SourceId).To(Equal("app-id"))
		Expect(env.InstanceId).To(Equal("source-instance"))

		ts := time.Unix(0, env.Timestamp)
		Expect(ts).Should(BeTemporally("~", time.Now(), time.Second))
		log := env.GetLog()
		Expect(log).NotTo(BeNil())
		Expect(log.Payload).To(Equal([]byte("message")))
		Expect(log.Type).To(Equal(loggregator_v2.Log_OUT))
	})

	It("sends logs", func() {
		client.EmitLog(
			"message",
			loggregator.WithSourceInfo("source-id", "source-type", "source-instance"),
			loggregator.WithStdout(),
		)
		env, err := getEnvelopeAt(server.receivers, 0)
		Expect(err).NotTo(HaveOccurred())

		Expect(env.SourceId).To(Equal("source-id"))
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
			loggregator.WithAppInfo("app-id", "source-type", "source-instance"),
		)

		env, err := getEnvelopeAt(server.receivers, 0)
		Expect(err).NotTo(HaveOccurred())

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
			loggregator.WithGaugeValue("name-a", 1, "unit-a"),
			loggregator.WithGaugeValue("name-b", 2, "unit-b"),
			loggregator.WithEnvelopeTags(map[string]string{"some-tag": "some-tag-value"}),
			loggregator.WithGaugeAppInfo("app-id", 123),
		)

		env, err := getEnvelopeAt(server.receivers, 0)
		Expect(err).NotTo(HaveOccurred())

		ts := time.Unix(0, env.Timestamp)
		Expect(ts).Should(BeTemporally("~", time.Now(), time.Second))
		metrics := env.GetGauge()
		Expect(metrics).NotTo(BeNil())
		Expect(env.SourceId).To(Equal("app-id"))
		Expect(env.InstanceId).To(Equal("123"))
		Expect(metrics.GetMetrics()).To(HaveLen(2))
		Expect(metrics.GetMetrics()["name-a"].Value).To(Equal(1.0))
		Expect(metrics.GetMetrics()["name-b"].Value).To(Equal(2.0))
		Expect(env.Tags["some-tag"]).To(Equal("some-tag-value"))
	})

	It("sends gauge metrics", func() {
		client.EmitGauge(
			loggregator.WithGaugeValue("name-a", 1, "unit-a"),
			loggregator.WithGaugeValue("name-b", 2, "unit-b"),
			loggregator.WithEnvelopeTags(map[string]string{"some-tag": "some-tag-value"}),
			loggregator.WithGaugeSourceInfo("source-id", "instance-id"),
		)

		env, err := getEnvelopeAt(server.receivers, 0)
		Expect(err).NotTo(HaveOccurred())

		ts := time.Unix(0, env.Timestamp)
		Expect(ts).Should(BeTemporally("~", time.Now(), time.Second))
		metrics := env.GetGauge()
		Expect(metrics).NotTo(BeNil())
		Expect(env.SourceId).To(Equal("source-id"))
		Expect(env.InstanceId).To(Equal("instance-id"))
		Expect(metrics.GetMetrics()).To(HaveLen(2))
		Expect(metrics.GetMetrics()["name-a"].Value).To(Equal(1.0))
		Expect(metrics.GetMetrics()["name-b"].Value).To(Equal(2.0))
		Expect(env.Tags["some-tag"]).To(Equal("some-tag-value"))
	})

	It("sends timers", func() {
		stopTime := time.Now()
		startTime := stopTime.Add(-time.Minute)

		client.EmitTimer("http", startTime, stopTime,
			loggregator.WithEnvelopeTags(map[string]string{"some-tag": "some-tag-value"}),
			loggregator.WithTimerSourceInfo("source-id", "instance-id"),
		)

		env, err := getEnvelopeAt(server.receivers, 0)
		Expect(err).ToNot(HaveOccurred())

		Expect(env.GetSourceId()).To(Equal("source-id"))
		Expect(env.GetInstanceId()).To(Equal("instance-id"))
		Expect(env.Tags["some-tag"]).To(Equal("some-tag-value"))

		timer := env.GetTimer()
		Expect(timer).ToNot(BeNil())
		Expect(timer.GetName()).To(Equal("http"))
		Expect(timer.GetStart()).To(Equal(startTime.UnixNano()))
		Expect(timer.GetStop()).To(Equal(stopTime.UnixNano()))
	})

	It("works with the runtime emitter", func() {
		// This test is to ensure that the v2 client satisfies the
		// runtimeemitter.Sender interface. If it does not satisfy the
		// runtimeemitter.Sender interface this test will force a compile time
		// error.
		runtimeemitter.New(client)
	})

	DescribeTable("emitting different envelope types", func(emit func()) {
		emit()

		env, err := getEnvelopeAt(server.receivers, 0)
		Expect(err).NotTo(HaveOccurred())

		Expect(env.Tags["string"]).To(Equal("client-string-tag"), "The client tag for string was not set properly")
		Expect(env.Tags["envelope-string"]).To(Equal("envelope-string-tag"), "The envelope tag for string was not set properly")
	},
		Entry("logs", func() {
			client.EmitLog(
				"message",
				loggregator.WithEnvelopeTag("envelope-string", "envelope-string-tag"),
			)
		}),
		Entry("gauge", func() {
			client.EmitGauge(
				loggregator.WithGaugeValue("gauge-name", 123.4, "some-unit"),
				loggregator.WithEnvelopeTag("envelope-string", "envelope-string-tag"),
			)
		}),
		Entry("counter", func() {
			client.EmitCounter(
				"foo",
				loggregator.WithEnvelopeTag("envelope-string", "envelope-string-tag"),
			)
		}),
	)

	It("sets the counter's delta to the given value", func() {
		e := &loggregator_v2.Envelope{
			Message: &loggregator_v2.Envelope_Counter{
				Counter: &loggregator_v2.Counter{},
			},
		}
		loggregator.WithDelta(99)(e)
		Expect(e.GetCounter().GetDelta()).To(Equal(uint64(99)))
	})

	It("sets the app info for a counter", func() {
		e := &loggregator_v2.Envelope{
			Message: &loggregator_v2.Envelope_Counter{
				Counter: &loggregator_v2.Counter{},
			},
		}
		loggregator.WithCounterAppInfo("some-guid", 101)(e)
		Expect(e.GetSourceId()).To(Equal("some-guid"))
		Expect(e.GetInstanceId()).To(Equal("101"))
	})

	It("sets the source info for a counter", func() {
		e := &loggregator_v2.Envelope{
			Message: &loggregator_v2.Envelope_Counter{
				Counter: &loggregator_v2.Counter{},
			},
		}
		loggregator.WithCounterSourceInfo("source-id", "instance-id")(e)
		Expect(e.GetSourceId()).To(Equal("source-id"))
		Expect(e.GetInstanceId()).To(Equal("instance-id"))
	})

	It("sets the title and body of an event envelope", func() {
		Eventually(func() error {
			return client.EmitEvent(
				context.Background(),
				"some-title",
				"some-body",
			)
		}).Should(Succeed())

		var envelopeBatch *loggregator_v2.EnvelopeBatch
		Eventually(server.sendReceiver).Should(Receive(&envelopeBatch))

		env := envelopeBatch.GetBatch()[0]
		Expect(env.GetEvent()).ToNot(BeNil())
		Expect(env.GetEvent().GetTitle()).To(Equal("some-title"))
		Expect(env.GetEvent().GetBody()).To(Equal("some-body"))
	})

	It("flushes current batch and sends", func() {
		client := buildIngressClient(server.addr, time.Hour)

		// Ensure client/server are ready
		Eventually(func() error {
			return client.EmitEvent(
				context.Background(),
				"some-title",
				"some-body",
			)
		}).Should(Succeed())

		client.EmitLog("message")
		err := client.CloseSend()
		Expect(err).ToNot(HaveOccurred())

		_, err = getEnvelopeAt(server.receivers, 0)
		Expect(err).ToNot(HaveOccurred())
	})

	It("does not block on an empty buffer", func(done Done) {
		defer close(done)

		err := client.CloseSend()
		Expect(err).ToNot(HaveOccurred())
	})
})

func getEnvelopeAt(receivers chan loggregator_v2.Ingress_BatchSenderServer, idx int) (*loggregator_v2.Envelope, error) {
	var recv loggregator_v2.Ingress_BatchSenderServer
	Eventually(receivers, 10).Should(Receive(&recv))

	envBatch, err := recv.Recv()
	if err != nil {
		return nil, err
	}

	if len(envBatch.Batch) < 1 {
		return nil, errors.New("no envelopes")
	}

	return envBatch.Batch[idx], nil
}

func buildIngressClient(serverAddr string, flushInterval time.Duration) *loggregator.IngressClient {
	tlsConfig, err := loggregator.NewIngressTLSConfig(
		fixture("CA.crt"),
		fixture("client.crt"),
		fixture("client.key"),
	)
	if err != nil {
		panic(err)
	}

	client, err := loggregator.NewIngressClient(
		tlsConfig,
		loggregator.WithAddr(serverAddr),
		loggregator.WithBatchFlushInterval(flushInterval),
		loggregator.WithTag("string", "client-string-tag"),
	)
	if err != nil {
		panic(err)
	}

	return client
}
