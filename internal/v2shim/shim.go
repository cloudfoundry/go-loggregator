package v2shim

import (
	"time"

	"code.cloudfoundry.org/go-loggregator/v2"
	"github.com/cloudfoundry/sonde-go/events"
)

type client struct {
	client *v2.Client
}

func NewClient(c *v2.Client) client {
	return client{client: c}
}

func (c client) SendDuration(name string, value time.Duration) error {
	c.client.SendDuration(name, value)

	return nil
}

func (c client) SendMebiBytes(name string, value int) error {
	c.client.SendMebiBytes(name, value)

	return nil
}

func (c client) SendMetric(name string, value int) error {
	c.client.SendMetric(name, value)

	return nil
}

func (c client) SendBytesPerSecond(name string, value float64) error {
	c.client.SendBytesPerSecond(name, value)

	return nil
}

func (c client) SendRequestsPerSecond(name string, value float64) error {
	c.client.SendRequestsPerSecond(name, value)

	return nil
}

func (c client) IncrementCounter(name string) error {
	c.client.EmitCounter(name)

	return nil
}

func (c client) SendAppLog(appID, message, sourceType, sourceInstance string) error {
	c.client.EmitLog(
		message,
		v2.WithAppInfo(appID, sourceType, sourceInstance),
		v2.WithStdout(),
	)
	return nil
}

func (c client) SendAppErrorLog(appID, message, sourceType, sourceInstance string) error {
	c.client.EmitLog(
		message,
		v2.WithAppInfo(appID, sourceType, sourceInstance),
	)
	return nil
}

func (c client) SendAppMetrics(metrics *events.ContainerMetric) error {
	c.client.SendAppMetrics(metrics)

	return nil
}
