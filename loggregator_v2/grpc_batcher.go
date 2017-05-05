package loggregator_v2

import "time"

type grpcBatcher struct {
	client  *grpcClient
	metrics map[string]*GaugeValue
}

func (b *grpcBatcher) Send() error {
	return b.client.send(&Envelope{
		Timestamp: time.Now().UnixNano(),
		Message: &Envelope_Gauge{
			Gauge: &Gauge{
				Metrics: b.metrics,
			},
		},
	})
}

func (b *grpcBatcher) SendDuration(name string, duration time.Duration) error {
	b.metrics[name] = &GaugeValue{
		Unit:  "nanos",
		Value: float64(duration),
	}
	return nil
}

func (b *grpcBatcher) SendMebiBytes(name string, mebibytes int) error {
	b.metrics[name] = &GaugeValue{
		Unit:  "MiB",
		Value: float64(mebibytes),
	}
	return nil
}

func (b *grpcBatcher) SendMetric(name string, value int) error {
	b.metrics[name] = &GaugeValue{
		Unit:  "Metric",
		Value: float64(value),
	}
	return nil
}

func (b *grpcBatcher) SendBytesPerSecond(name string, value float64) error {
	b.metrics[name] = &GaugeValue{
		Unit:  "B/s",
		Value: value,
	}
	return nil
}

func (b *grpcBatcher) SendRequestsPerSecond(name string, value float64) error {
	b.metrics[name] = &GaugeValue{
		Unit:  "Req/s",
		Value: value,
	}
	return nil
}
