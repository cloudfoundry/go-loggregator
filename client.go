package loggregator

import (
	"time"

	"code.cloudfoundry.org/go-loggregator/v1"
	"code.cloudfoundry.org/go-loggregator/v2"

	"github.com/cloudfoundry/sonde-go/events"
)

type Client interface {
	SendDuration(name string, value time.Duration)
	SendMebiBytes(name string, value int)
	SendMetric(name string, value int)
	SendBytesPerSecond(name string, value float64)
	SendRequestsPerSecond(name string, value float64)
	IncrementCounter(name string)
	SendAppLog(appID, message, sourceType, sourceInstance string)
	SendAppErrorLog(appID, message, sourceType, sourceInstance string)
	SendAppMetrics(metrics *events.ContainerMetric)
}

type Config struct {
	UseV2API           bool
	APIPort            int
	CACertPath         string
	CertPath           string
	KeyPath            string
	JobDeployment      string
	JobName            string
	JobIndex           string
	JobIP              string
	JobOrigin          string
	BatchMaxSize       uint
	BatchFlushInterval time.Duration
}

// NewClient creates a connection to the Loggregator API.
// Users can opt-in to using the v2 API through configuration.
// If an opt-in feature is not required, the v1 and v2 clients are both
// available for public use as well.
func NewClient(config Config) (Client, error) {
	if config.UseV2API {
		return newV2Client(config)
	}

	return v1.NewClient()
}

func newV2Client(config Config) (Client, error) {
	tlsConfig, err := v2.NewTLSConfig(
		config.CACertPath,
		config.CertPath,
		config.KeyPath,
	)
	if err != nil {
		return nil, err
	}

	var opts []v2.Option

	if config.BatchMaxSize != 0 {
		opts = append(opts, v2.WithBatchMaxSize(config.BatchMaxSize))
	}

	if config.BatchFlushInterval != time.Duration(0) {
		opts = append(opts, v2.WithBatchFlushInterval(config.BatchFlushInterval))
	}

	return v2.NewClient(tlsConfig, config.APIPort, opts...)
}
