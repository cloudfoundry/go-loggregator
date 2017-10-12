package loggregator_test

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
)

type testIngressServer struct {
	receivers  chan loggregator_v2.Ingress_BatchSenderServer
	addr       string
	tlsConfig  *tls.Config
	grpcServer *grpc.Server
	grpc.Stream
}

func newTestIngressServer(serverCert, serverKey, caCert string) (*testIngressServer, error) {
	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequestClientCert,
		InsecureSkipVerify: false,
	}
	caCertBytes, err := ioutil.ReadFile(caCert)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertBytes)
	tlsConfig.RootCAs = caCertPool

	return &testIngressServer{
		tlsConfig: tlsConfig,
		receivers: make(chan loggregator_v2.Ingress_BatchSenderServer),
		addr:      "localhost:0",
	}, nil
}

func (*testIngressServer) Sender(srv loggregator_v2.Ingress_SenderServer) error {
	return nil
}

func (t *testIngressServer) BatchSender(srv loggregator_v2.Ingress_BatchSenderServer) error {
	t.receivers <- srv

	<-srv.Context().Done()

	return nil
}

func (t *testIngressServer) start() error {
	listener, err := net.Listen("tcp4", t.addr)
	if err != nil {
		return err
	}
	t.addr = listener.Addr().String()

	var opts []grpc.ServerOption
	if t.tlsConfig != nil {
		opts = append(opts, grpc.Creds(credentials.NewTLS(t.tlsConfig)))
	}
	t.grpcServer = grpc.NewServer(opts...)

	loggregator_v2.RegisterIngressServer(t.grpcServer, t)

	go t.grpcServer.Serve(listener)

	return nil
}

func (t *testIngressServer) stop() {
	t.grpcServer.Stop()
}

type testEgressServer struct {
	addr_      string
	cn         string
	tlsConfig  *tls.Config
	grpcServer *grpc.Server
	grpc.Stream
}

type egressServerOption func(*testEgressServer)

func withCN(cn string) egressServerOption {
	return func(s *testEgressServer) {
		s.cn = cn
	}
}

func withAddr(addr string) egressServerOption {
	return func(s *testEgressServer) {
		s.addr_ = addr
	}
}

func newTestEgressServer(serverCert, serverKey, caCert string, opts ...egressServerOption) (*testEgressServer, error) {
	s := &testEgressServer{
		addr_: "localhost:0",
	}

	for _, o := range opts {
		o(s)
	}

	cert, err := tls.LoadX509KeyPair(serverCert, serverKey)
	if err != nil {
		return nil, err
	}

	s.tlsConfig = &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequestClientCert,
		InsecureSkipVerify: false,
		ServerName:         s.cn,
	}
	caCertBytes, err := ioutil.ReadFile(caCert)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertBytes)
	s.tlsConfig.RootCAs = caCertPool

	return s, nil
}

func (t *testEgressServer) addr() string {
	return t.addr_
}

func (t *testEgressServer) Receiver(*loggregator_v2.EgressRequest, loggregator_v2.Egress_ReceiverServer) error {
	return nil
}

func (t *testEgressServer) BatchedReceiver(*loggregator_v2.EgressBatchRequest, loggregator_v2.Egress_BatchedReceiverServer) error {
	return nil
}

func (t *testEgressServer) start() error {
	listener, err := net.Listen("tcp4", t.addr_)
	if err != nil {
		return err
	}
	t.addr_ = listener.Addr().String()

	var opts []grpc.ServerOption
	if t.tlsConfig != nil {
		opts = append(opts, grpc.Creds(credentials.NewTLS(t.tlsConfig)))
	}
	t.grpcServer = grpc.NewServer(opts...)

	loggregator_v2.RegisterEgressServer(t.grpcServer, t)

	go t.grpcServer.Serve(listener)

	return nil
}

func (t *testEgressServer) stop() {
	t.grpcServer.Stop()
}
