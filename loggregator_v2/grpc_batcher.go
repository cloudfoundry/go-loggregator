package loggregator_v2

import "time"

type grpcBatcher struct {
	c       *grpcClient
	metrics map[string]*GaugeValue
}

func (gb *grpcBatcher) Send() error {
	return gb.c.send(&Envelope{
		Timestamp: time.Now().UnixNano(),
		Message: &Envelope_Gauge{
			Gauge: &Gauge{
				Metrics: gb.metrics,
			},
		},
	})
}

func (c *grpcBatcher) SendDuration(name string, duration time.Duration) error {
	c.metrics[name] = &GaugeValue{
		Unit:  "nanos",
		Value: float64(duration),
	}
	return nil
}

func (c *grpcBatcher) SendMebiBytes(name string, mebibytes int) error {
	c.metrics[name] = &GaugeValue{
		Unit:  "MiB",
		Value: float64(mebibytes),
	}
	return nil
}

func (c *grpcBatcher) SendMetric(name string, value int) error {
	c.metrics[name] = &GaugeValue{
		Unit:  "Metric",
		Value: float64(value),
	}
	return nil
}

func (c *grpcBatcher) SendBytesPerSecond(name string, value float64) error {
	c.metrics[name] = &GaugeValue{
		Unit:  "B/s",
		Value: value,
	}
	return nil
}

func (c *grpcBatcher) SendRequestsPerSecond(name string, value float64) error {
	c.metrics[name] = &GaugeValue{
		Unit:  "Req/s",
		Value: value,
	}
	return nil
}
