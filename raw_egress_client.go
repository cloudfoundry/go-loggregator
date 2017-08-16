package loggregator

import (
	"crypto/tls"
	"io"

	"golang.org/x/net/context"

	"code.cloudfoundry.org/go-loggregator/rpc/loggregator_v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// RawEgressClient wraps the gRPC RawEgressClient for convenience.
type RawEgressClient struct {
	c loggregator_v2.EgressClient
}

// NewRawEgressClient creates a new EgressClient for the given addr and TLS
// configuration.
func NewRawEgressClient(addr string, c *tls.Config) (*RawEgressClient, io.Closer, error) {
	conn, err := grpc.Dial(addr,
		grpc.WithTransportCredentials(credentials.NewTLS(c)),
	)
	if err != nil {
		return nil, nil, err
	}

	return &RawEgressClient{c: loggregator_v2.NewEgressClient(conn)}, conn, nil
}

// Receiver wraps the created EgressClient's Receiver method.
func (c *RawEgressClient) Receiver(
	ctx context.Context,
	in *loggregator_v2.EgressRequest,
) (loggregator_v2.Egress_ReceiverClient, error) {

	return c.c.Receiver(ctx, in)
}

func (c *RawEgressClient) BatchReceiver(
	ctx context.Context,
	in *loggregator_v2.EgressBatchRequest,
) (loggregator_v2.Egress_BatchedReceiverClient, error) {

	return c.c.BatchedReceiver(ctx, in)
}
