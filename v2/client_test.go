package v2_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/go-loggregator/internal/loggregator_v2"
	"code.cloudfoundry.org/go-loggregator/v2"

	"github.com/cloudfoundry/sonde-go/events"
	"github.com/gogo/protobuf/proto"
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
		}

		batch, err := getBatch(receivers)
		Expect(err).ToNot(HaveOccurred())

		Expect(len(batch.Batch)).To(BeNumerically(">", 1))
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
		metric := events.ContainerMetric{
			ApplicationId:    proto.String("app-id"),
			CpuPercentage:    proto.Float64(10.0),
			MemoryBytes:      proto.Uint64(10),
			DiskBytes:        proto.Uint64(10),
			MemoryBytesQuota: proto.Uint64(20),
			DiskBytesQuota:   proto.Uint64(20),
			InstanceIndex:    proto.Int32(5),
		}

		client.SendAppMetrics(&metric)

		env, err := getEnvelopeAt(receivers, 0)
		Expect(err).NotTo(HaveOccurred())

		ts := time.Unix(0, env.Timestamp)
		Expect(ts).Should(BeTemporally("~", time.Now(), time.Second))
		metrics := env.GetGauge()
		Expect(metrics).NotTo(BeNil())
		Expect(env.SourceId).To(Equal("app-id"))
		Expect(metrics.GetMetrics()).To(HaveLen(6))
		Expect(metrics.GetMetrics()["instance_index"].Value).To(Equal(5.0))
		Expect(metrics.GetMetrics()["cpu"].Value).To(Equal(10.0))
		Expect(metrics.GetMetrics()["memory"].Value).To(Equal(10.0))
		Expect(metrics.GetMetrics()["disk"].Value).To(Equal(10.0))
		Expect(metrics.GetMetrics()["memory_quota"].Value).To(Equal(20.0))
		Expect(metrics.GetMetrics()["disk_quota"].Value).To(Equal(20.0))
	})

	Context("when component metrics are emitted", func() {
		It("sends duration info", func() {
			client.SendDuration("test-name", 1*time.Nanosecond)

			env, err := getEnvelopeAt(receivers, 0)
			Expect(err).NotTo(HaveOccurred())

			ts := time.Unix(0, env.Timestamp)
			Expect(ts).Should(BeTemporally("~", time.Now(), time.Second))
			message := env.GetGauge()
			Expect(message).NotTo(BeNil())
			Expect(message.GetMetrics()["test-name"].Value).To(Equal(float64(1)))
			Expect(message.GetMetrics()["test-name"].Unit).To(Equal("nanos"))
		})

		It("sends mebibytes info", func() {
			client.SendMebiBytes("test-name", 10)

			env, err := getEnvelopeAt(receivers, 0)
			Expect(err).NotTo(HaveOccurred())

			ts := time.Unix(0, env.Timestamp)
			Expect(ts).Should(BeTemporally("~", time.Now(), time.Second))
			message := env.GetGauge()
			Expect(message).NotTo(BeNil())
			Expect(message.GetMetrics()["test-name"].Value).To(Equal(float64(10)))
			Expect(message.GetMetrics()["test-name"].Unit).To(Equal("MiB"))
		})

		It("sends metrics info", func() {
			client.SendMetric("test-name", 11)

			env, err := getEnvelopeAt(receivers, 0)
			Expect(err).NotTo(HaveOccurred())

			ts := time.Unix(0, env.Timestamp)
			Expect(ts).Should(BeTemporally("~", time.Now(), time.Second))
			message := env.GetGauge()
			Expect(message).NotTo(BeNil())
			Expect(message.GetMetrics()["test-name"].Value).To(Equal(float64(11)))
			Expect(message.GetMetrics()["test-name"].Unit).To(Equal("Metric"))
		})

		It("sends requests per second info", func() {
			client.SendRequestsPerSecond("test-name", 11)

			env, err := getEnvelopeAt(receivers, 0)
			Expect(err).NotTo(HaveOccurred())

			ts := time.Unix(0, env.Timestamp)
			Expect(ts).Should(BeTemporally("~", time.Now(), time.Second))
			message := env.GetGauge()
			Expect(message).NotTo(BeNil())
			Expect(message.GetMetrics()["test-name"].Value).To(Equal(float64(11)))
		})

		It("sends bytes per second info", func() {
			client.SendBytesPerSecond("test-name", 10)

			env, err := getEnvelopeAt(receivers, 0)
			Expect(err).NotTo(HaveOccurred())

			ts := time.Unix(0, env.Timestamp)
			Expect(ts).Should(BeTemporally("~", time.Now(), time.Second))
			message := env.GetGauge()
			Expect(message).NotTo(BeNil())
			Expect(message.GetMetrics()["test-name"].Value).To(Equal(float64(10)))
			Expect(message.GetMetrics()["test-name"].Unit).To(Equal("B/s"))
		})

		It("increments counter", func() {
			client.IncrementCounter("test-name")

			env, err := getEnvelopeAt(receivers, 0)
			Expect(err).NotTo(HaveOccurred())

			ts := time.Unix(0, env.Timestamp)
			Expect(ts).Should(BeTemporally("~", time.Now(), time.Second))
			message := env.GetCounter()
			Expect(message).NotTo(BeNil())
			Expect(message.Name).To(Equal("test-name"))
			Expect(message.GetDelta()).To(Equal(uint64(1)))
		})
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

		closeCh := make(chan struct{})
		go func() {
			for {
				select {
				case <-closeCh:
					break
				default:
					client.EmitLog(
						"message",
						v2.WithAppInfo("app-id", "source-type", "source-instance"),
					)

					time.Sleep(50 * time.Millisecond)
				}
			}
		}()
		defer close(closeCh)

		envBatch, err = getBatch(receivers)
		Expect(err).NotTo(HaveOccurred())
		Expect(envBatch.Batch).ToNot(BeEmpty())
	})
})

func getBatch(receivers chan loggregator_v2.Ingress_BatchSenderServer) (*loggregator_v2.EnvelopeBatch, error) {
	var recv loggregator_v2.Ingress_BatchSenderServer
	Eventually(receivers, 3).Should(Receive(&recv))

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
