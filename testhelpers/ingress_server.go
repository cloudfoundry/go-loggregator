package testhelpers

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"code.cloudfoundry.org/go-loggregator/testhelpers/fakes"
	"code.cloudfoundry.org/localip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type TestIngressServer struct {
	receivers  chan loggregator_v2.Ingress_BatchSenderServer
	port       int
	tlsConfig  *tls.Config
	grpcServer *grpc.Server
}

func NewTestIngressServer(serverCert, serverKey, caCert string) (*TestIngressServer, error) {
	port, err := localip.LocalPort()
	if err != nil {
		return nil, err
	}

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

	return &TestIngressServer{
		tlsConfig: tlsConfig,
		receivers: make(chan loggregator_v2.Ingress_BatchSenderServer),
		port:      int(port),
	}, nil
}

func (t *TestIngressServer) Port() int {
	return t.port
}

func (t *TestIngressServer) Receivers() chan loggregator_v2.Ingress_BatchSenderServer {
	return t.receivers
}

func (t *TestIngressServer) Start() error {
	listener, err := net.Listen("tcp4", fmt.Sprintf("localhost:%d", t.port))
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	if t.tlsConfig != nil {
		opts = append(opts, grpc.Creds(credentials.NewTLS(t.tlsConfig)))
	}
	t.grpcServer = grpc.NewServer(opts...)

	senderServer := &fakes.FakeIngressServer{}
	senderServer.BatchSenderStub = func(recv loggregator_v2.Ingress_BatchSenderServer) error {
		t.receivers <- recv
		return nil
	}
	loggregator_v2.RegisterIngressServer(t.grpcServer, senderServer)

	go t.grpcServer.Serve(listener)

	return nil
}

func (t *TestIngressServer) Stop() {
	t.grpcServer.Stop()
}
