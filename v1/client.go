package v1

import (
	"time"

	"github.com/cloudfoundry/dropsonde/logs"
	"github.com/cloudfoundry/dropsonde/metrics"
	"github.com/cloudfoundry/sonde-go/events"
)

func NewClient() (*Client, error) {
	return &Client{}, nil
}

type Client struct{}

func (c *Client) Send() error {
	return nil
}

func (c *Client) IncrementCounter(name string) {
	metrics.IncrementCounter(name)
}
func (c *Client) SendAppLog(appID, message, sourceType, sourceInstance string) {
	logs.SendAppLog(appID, message, sourceType, sourceInstance)
}

func (c *Client) SendAppErrorLog(appID, message, sourceType, sourceInstance string) {
	logs.SendAppErrorLog(appID, message, sourceType, sourceInstance)
}

func (c *Client) SendAppMetrics(m *events.ContainerMetric) {
	metrics.Send(m)
}

func (c *Client) SendDuration(name string, duration time.Duration) {
	c.sendComponentMetric(name, float64(duration), "nanos")
}

func (c *Client) SendMebiBytes(name string, mebibytes int) {
	c.sendComponentMetric(name, float64(mebibytes), "MiB")
}

func (c *Client) SendMetric(name string, value int) {
	c.sendComponentMetric(name, float64(value), "Metric")
}

func (c *Client) SendBytesPerSecond(name string, value float64) {
	c.sendComponentMetric(name, value, "B/s")
}

func (c *Client) SendRequestsPerSecond(name string, value float64) {
	c.sendComponentMetric(name, value, "Req/s")
}

func (c *Client) sendComponentMetric(name string, value float64, unit string) {
	metrics.SendValue(name, value, unit)
}
