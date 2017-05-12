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
	c.client.EmitGauge(
		v2.WithGaugeValue(name, float64(value), "nanos"),
	)
	return nil
}

func (c client) SendMebiBytes(name string, value int) error {
	c.client.EmitGauge(
		v2.WithGaugeValue(name, float64(value), "MiB"),
	)
	return nil
}

func (c client) SendMetric(name string, value int) error {
	c.client.EmitGauge(
		v2.WithGaugeValue(name, float64(value), "Metric"),
	)

	return nil
}

func (c client) SendBytesPerSecond(name string, value float64) error {
	c.client.EmitGauge(
		v2.WithGaugeValue(name, value, "B/s"),
	)
	return nil
}

func (c client) SendRequestsPerSecond(name string, value float64) error {
	c.client.EmitGauge(
		v2.WithGaugeValue(name, value, "Req/s"),
	)
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

func (c client) SendAppMetrics(m *events.ContainerMetric) error {
	c.client.EmitGauge(
		v2.WithGaugeValue("instance_index", float64(m.GetInstanceIndex()), ""),
		v2.WithGaugeValue("cpu", m.GetCpuPercentage(), "percentage"),
		v2.WithGaugeValue("memory", float64(m.GetMemoryBytes()), "bytes"),
		v2.WithGaugeValue("disk", float64(m.GetDiskBytes()), "bytes"),
		v2.WithGaugeValue("memory_quota", float64(m.GetMemoryBytesQuota()), "bytes"),
		v2.WithGaugeValue("disk_quota", float64(m.GetDiskBytesQuota()), "bytes"),
		v2.WithGaugeAppInfo(m.GetApplicationId()),
	)

	return nil
}
