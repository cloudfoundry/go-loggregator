package v1

import (
	"time"

	"github.com/cloudfoundry/dropsonde/logs"
	"github.com/cloudfoundry/dropsonde/metrics"
	"github.com/cloudfoundry/sonde-go/events"
)

func NewClient() (*dropsondeClient, error) {
	return &dropsondeClient{}, nil
}

type dropsondeClient struct{}

func (c *dropsondeClient) Send() error {
	return nil
}

func (c *dropsondeClient) IncrementCounter(name string) {
	metrics.IncrementCounter(name)
}
func (c *dropsondeClient) SendAppLog(appID, message, sourceType, sourceInstance string) {
	logs.SendAppLog(appID, message, sourceType, sourceInstance)
}

func (c *dropsondeClient) SendAppErrorLog(appID, message, sourceType, sourceInstance string) {
	logs.SendAppErrorLog(appID, message, sourceType, sourceInstance)
}

func (c *dropsondeClient) SendAppMetrics(m *events.ContainerMetric) {
	metrics.Send(m)
}

func (c *dropsondeClient) SendDuration(name string, duration time.Duration) {
	c.sendComponentMetric(name, float64(duration), "nanos")
}

func (c *dropsondeClient) SendMebiBytes(name string, mebibytes int) {
	c.sendComponentMetric(name, float64(mebibytes), "MiB")
}

func (c *dropsondeClient) SendMetric(name string, value int) {
	c.sendComponentMetric(name, float64(value), "Metric")
}

func (c *dropsondeClient) SendBytesPerSecond(name string, value float64) {
	c.sendComponentMetric(name, value, "B/s")
}

func (c *dropsondeClient) SendRequestsPerSecond(name string, value float64) {
	c.sendComponentMetric(name, value, "Req/s")
}

func (c *dropsondeClient) sendComponentMetric(name string, value float64, unit string) {
	metrics.SendValue(name, value, unit)
}
