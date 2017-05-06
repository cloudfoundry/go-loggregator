package loggregator_v2_test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"

	"code.cloudfoundry.org/go-loggregator/loggregator_v2"
	"code.cloudfoundry.org/go-loggregator/loggregator_v2/fakes"
	"code.cloudfoundry.org/localip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type TestServer struct {
	serverCert string
	serverKey  string
	caCert     string

	receivers chan loggregator_v2.Ingress_BatchSenderServer
	port      int

	grpcServer *grpc.Server
}

func NewTestServer(serverCert, serverKey, caCert string) (*TestServer, error) {
	port, err := localip.LocalPort()
	if err != nil {
		return nil, err
	}

	return &TestServer{
		serverCert: serverCert,
		serverKey:  serverKey,
		caCert:     caCert,
		receivers:  make(chan loggregator_v2.Ingress_BatchSenderServer),
		port:       int(port),
	}, nil
}

func (t *TestServer) Port() int {
	return t.port
}

func (t *TestServer) Receivers() chan loggregator_v2.Ingress_BatchSenderServer {
	return t.receivers
}

func (t *TestServer) Start() error {
	cert, err := tls.LoadX509KeyPair(t.serverCert, t.serverKey)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequestClientCert,
		InsecureSkipVerify: false,
	}
	caCertBytes, err := ioutil.ReadFile(t.caCert)
	if err != nil {
		return err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertBytes)
	tlsConfig.RootCAs = caCertPool
	t.grpcServer = grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))

	senderServer := &fakes.FakeIngressServer{}
	senderServer.BatchSenderStub = func(recv loggregator_v2.Ingress_BatchSenderServer) error {
		t.receivers <- recv
		return nil
	}
	loggregator_v2.RegisterIngressServer(t.grpcServer, senderServer)

	listener, err := net.Listen("tcp4", fmt.Sprintf("localhost:%d", t.port))
	if err != nil {
		return err
	}

	go func() {
		t.grpcServer.Serve(listener)
	}()

	return nil
}

func (t *TestServer) Stop() {
	t.grpcServer.Stop()
}
