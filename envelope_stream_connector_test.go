package loggregator_test

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"

	"code.cloudfoundry.org/go-loggregator/v8"
	"code.cloudfoundry.org/go-loggregator/v8/rpc/loggregator_v2"
	"code.cloudfoundry.org/tlsconfig"
	"github.com/gogo/protobuf/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Connector", func() {
	It("initiates a connection to receive envelopes", func() {
		producer, err := newFakeEventProducer()
		Expect(err).NotTo(HaveOccurred())
		producer.start()
		defer producer.stop()
		tlsConf, err := loggregator.NewIngressTLSConfig(
			fixture("CA.crt"),
			fixture("server.crt"),
			fixture("server.key"),
		)
		Expect(err).NotTo(HaveOccurred())

		addr := producer.addr
		req := &loggregator_v2.EgressBatchRequest{ShardId: "some-id"}
		c := loggregator.NewEnvelopeStreamConnector(
			addr,
			tlsConf,
		)

		rx := c.Stream(context.Background(), req)

		Expect(len(rx())).NotTo(BeZero())
		Expect(proto.Equal(producer.actualReq(), req)).To(BeTrue())
	})

	It("reconnects if the stream fails", func() {
		producer, err := newFakeEventProducer()
		Expect(err).NotTo(HaveOccurred())

		// Producer will grab a port on start. When the producer is restarted,
		// it will grab the same port.
		producer.start()

		tlsConf, err := loggregator.NewIngressTLSConfig(
			fixture("CA.crt"),
			fixture("server.crt"),
			fixture("server.key"),
		)
		Expect(err).NotTo(HaveOccurred())

		addr := producer.addr
		c := loggregator.NewEnvelopeStreamConnector(
			addr,
			tlsConf,
		)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		go func(ctx context.Context) {
			rx := c.Stream(ctx, &loggregator_v2.EgressBatchRequest{})
			for {
				if ctx.Err() != nil {
					return
				}
				rx()
			}
		}(ctx)

		Eventually(producer.connectionAttempts).Should(Equal(1))
		producer.stop()
		producer.start()
		defer producer.stop()

		// Reconnect after killing the server.
		Eventually(producer.connectionAttempts, 5).Should(Equal(2))

		// Ensure we don't create new connections when we don't need to.
		Consistently(producer.connectionAttempts).Should(Equal(2))
	})

	It("enables buffering", func() {
		producer, err := newFakeEventProducer()
		Expect(err).NotTo(HaveOccurred())

		// Producer will grab a port on start. When the producer is restarted,
		// it will grab the same port.
		producer.start()
		defer producer.stop()

		tlsConf, err := loggregator.NewIngressTLSConfig(
			fixture("CA.crt"),
			fixture("server.crt"),
			fixture("server.key"),
		)
		Expect(err).NotTo(HaveOccurred())

		var (
			mu     sync.Mutex
			missed int
		)
		addr := producer.addr
		c := loggregator.NewEnvelopeStreamConnector(
			addr,
			tlsConf,
			loggregator.WithEnvelopeStreamBuffer(5, func(m int) {
				mu.Lock()
				defer mu.Unlock()
				missed += m
			}),
		)
		rx := c.Stream(context.Background(), &loggregator_v2.EgressBatchRequest{})

		var count int
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		// Read to allow the diode to notice it dropped data
		go func(ctx context.Context) {
			for {
				select {
				case <-time.Tick(500 * time.Millisecond):
					// Do not invoke rx while mu is locked
					l := len(rx())
					mu.Lock()
					count += l
					mu.Unlock()
				case <-ctx.Done():
					return
				}
			}
		}(ctx)

		Eventually(func() int {
			mu.Lock()
			defer mu.Unlock()
			return missed
		}).ShouldNot(BeZero())

		mu.Lock()
		l := count
		mu.Unlock()
		Expect(l).ToNot(BeZero())
	})

	It("won't panic when context canceled", func() {
		producer, err := newFakeEventProducer()
		Expect(err).NotTo(HaveOccurred())

		tlsConf, err := loggregator.NewIngressTLSConfig(
			fixture("CA.crt"),
			fixture("server.crt"),
			fixture("server.key"),
		)
		Expect(err).NotTo(HaveOccurred())

		c := loggregator.NewEnvelopeStreamConnector(
			producer.addr,
			tlsConf,
			loggregator.WithEnvelopeStreamBuffer(5, func(m int) {}),
		)

		ctx, cancel := context.WithCancel(context.Background())
		rx := c.Stream(ctx, &loggregator_v2.EgressBatchRequest{})

		cancel()
		msg := rx()
		Expect(msg).To(BeNil())
	})
})

type fakeEventProducer struct {
	loggregator_v2.UnimplementedEgressServer

	server *grpc.Server
	addr   string

	mu                  sync.Mutex
	connectionAttempts_ int
	actualReq_          *loggregator_v2.EgressBatchRequest
}

func newFakeEventProducer() (*fakeEventProducer, error) {
	f := &fakeEventProducer{}

	return f, nil
}

func (f *fakeEventProducer) Receiver(
	*loggregator_v2.EgressRequest,
	loggregator_v2.Egress_ReceiverServer,
) error {

	return status.Errorf(codes.Unimplemented, "use BatchedReceiver instead")
}

func (f *fakeEventProducer) BatchedReceiver(
	req *loggregator_v2.EgressBatchRequest,
	srv loggregator_v2.Egress_BatchedReceiverServer,
) error {
	f.mu.Lock()
	f.connectionAttempts_++
	f.actualReq_ = req
	f.mu.Unlock()
	var i int
	for range time.Tick(10 * time.Millisecond) {
		err := srv.Send(&loggregator_v2.EnvelopeBatch{
			Batch: []*loggregator_v2.Envelope{
				{
					SourceId: fmt.Sprintf("envelope-%d", i),
					Message: &loggregator_v2.Envelope_Event{
						Event: &loggregator_v2.Event{
							Title: "event-name",
							Body:  "event-body",
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
		i++
	}
	return nil
}

func (f *fakeEventProducer) start() {
	addr := f.addr
	if addr == "" {
		addr = "127.0.0.1:0"
	}
	var lis net.Listener
	for i := 0; ; i++ {
		var err error
		lis, err = net.Listen("tcp", addr)
		if err != nil {
			// This can happen if the port is already in use...
			if i < 50 {
				log.Printf("failed to bind for fake producer. Trying again (%d/50)...: %s", i+1, err)
				time.Sleep(100 * time.Millisecond)
				continue
			}
			panic(err)
		}
		break
	}
	f.addr = lis.Addr().String()
	c, err := newServerMutualTLSConfig()
	if err != nil {
		panic(err)
	}
	opt := grpc.Creds(credentials.NewTLS(c))
	f.server = grpc.NewServer(opt)
	loggregator_v2.RegisterEgressServer(f.server, f)

	go f.listen(lis)
}

func (f *fakeEventProducer) listen(lis net.Listener) {
	_ = f.server.Serve(lis)
}

func (f *fakeEventProducer) stop() bool {
	if f.server == nil {
		return false
	}

	f.server.Stop()
	f.server = nil
	return true
}

func (f *fakeEventProducer) actualReq() *loggregator_v2.EgressBatchRequest {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.actualReq_
}

func (f *fakeEventProducer) connectionAttempts() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.connectionAttempts_
}

func newServerMutualTLSConfig() (*tls.Config, error) {
	certFile := fixture("server.crt")
	keyFile := fixture("server.key")
	caCertFile := fixture("CA.crt")

	return tlsconfig.Build(
		tlsconfig.WithInternalServiceDefaults(),
		tlsconfig.WithIdentityFromFile(certFile, keyFile),
	).Server(
		tlsconfig.WithClientAuthenticationFromFile(caCertFile),
	)
}
