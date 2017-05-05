package loggregator_v2_test

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"os"

	"code.cloudfoundry.org/go-loggregator/loggregator_v2"
	"code.cloudfoundry.org/go-loggregator/loggregator_v2/fakes"
	"code.cloudfoundry.org/localip"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GrpcRunner struct {
	serverCert     string
	serverKey      string
	caCert         string
	receivers      chan loggregator_v2.Ingress_SenderServer
	batchReceivers chan loggregator_v2.Ingress_BatchSenderServer
	port           int
}

func NewGRPCRunner(serverCert, serverKey, caCert string) (*GrpcRunner, error) {
	port, err := localip.LocalPort()
	if err != nil {
		return nil, err
	}

	return &GrpcRunner{
		serverCert:     serverCert,
		serverKey:      serverKey,
		caCert:         caCert,
		receivers:      make(chan loggregator_v2.Ingress_SenderServer),
		batchReceivers: make(chan loggregator_v2.Ingress_BatchSenderServer),
		port:           int(port),
	}, nil
}

func (grpcRunner *GrpcRunner) Port() int {
	return grpcRunner.port
}

func (grpcRunner *GrpcRunner) Receivers() chan loggregator_v2.Ingress_SenderServer {
	return grpcRunner.receivers
}

func (grpcRunner *GrpcRunner) BatchReceivers() chan loggregator_v2.Ingress_BatchSenderServer {
	return grpcRunner.batchReceivers
}

func (grpcRunner *GrpcRunner) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	cert, err := tls.LoadX509KeyPair(grpcRunner.serverCert, grpcRunner.serverKey)
	if err != nil {
		return err
	}

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		ClientAuth:         tls.RequestClientCert,
		InsecureSkipVerify: false,
	}
	caCertBytes, err := ioutil.ReadFile(grpcRunner.caCert)
	if err != nil {
		return err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCertBytes)
	tlsConfig.RootCAs = caCertPool
	server := grpc.NewServer(grpc.Creds(credentials.NewTLS(tlsConfig)))

	senderServer := &fakes.FakeIngressServer{}
	senderServer.SenderStub = func(recv loggregator_v2.Ingress_SenderServer) error {
		grpcRunner.receivers <- recv
		return nil
	}
	senderServer.BatchSenderStub = func(recv loggregator_v2.Ingress_BatchSenderServer) error {
		grpcRunner.batchReceivers <- recv
		return nil
	}
	loggregator_v2.RegisterIngressServer(server, senderServer)

	listener, err := net.Listen("tcp4", fmt.Sprintf("localhost:%d", grpcRunner.port))
	if err != nil {
		return err
	}

	errCh := make(chan error)
	go func() {
		errCh <- server.Serve(listener)
	}()
	close(ready)
	select {
	case <-signals:
		server.Stop()
		return nil
	case err := <-errCh:
		return err
	}
}
