// loggregator provides a top-level client for connecting to the loggregator v1
// and v2 API's.
//
// All members in the package here are deprecated and will be removed in the
// next major version of this library. Instead, see the v1 and v2 packages for
// the preferred way of connecting to the respective loggregator API.
package loggregator

import (
	"time"

	"code.cloudfoundry.org/go-loggregator/internal/v2shim"
	"code.cloudfoundry.org/go-loggregator/v1"
	"code.cloudfoundry.org/go-loggregator/v2"

	"github.com/cloudfoundry/sonde-go/events"
)

// Client is the shared contract between v1 and v2 clients.
//
// Deprecated: This interface will be removed in the next major version.
// Instead, use the v1 or v2 clients directly.
type Client interface {
	SendDuration(name string, value time.Duration) error
	SendMebiBytes(name string, value int) error
	SendMetric(name string, value int) error
	SendBytesPerSecond(name string, value float64) error
	SendRequestsPerSecond(name string, value float64) error
	IncrementCounter(name string) error
	SendAppLog(appID, message, sourceType, sourceInstance string) error
	SendAppErrorLog(appID, message, sourceType, sourceInstance string) error
	SendAppMetrics(metrics *events.ContainerMetric) error
}

// Config is the shared configuration between v1 and v2 clients.
//
// Deprecated: Config will be removed in the next major version.
// Instead, create a v1 or v2 client directly.
type Config struct {
	UseV2API      bool   `json:"loggregator_use_v2_api"`
	APIPort       int    `json:"loggregator_api_port"`
	CACertPath    string `json:"loggregator_ca_path"`
	CertPath      string `json:"loggregator_cert_path"`
	KeyPath       string `json:"loggregator_key_path"`
	JobDeployment string `json:"loggregator_job_deployment"`
	JobName       string `json:"loggregator_job_name"`
	JobIndex      string `json:"loggregator_job_index"`
	JobIP         string `json:"loggregator_job_ip"`
	JobOrigin     string `json:"loggregator_job_origin"`

	BatchMaxSize       uint
	BatchFlushInterval time.Duration
}

// NewClient returns a v1 or v2 client depending on the value of `UseV2API`
// from the config
//
// Deprecated: NewClient will be removed in the next major version.
// Instead, create a v1 or v2 client directly.
func NewClient(config Config) (Client, error) {
	if config.UseV2API {
		return NewV2Client(config)
	}

	return NewV1Client(config)
}

// NewV1Client creates a V1 connection to the Loggregator API.
//
// Deprecated: NewV1Client will be removed in the next major version.
// Instead, use v1.NewClient.
func NewV1Client(config Config) (Client, error) {
	return v1.NewClient()
}

// NewV2Client creates a V2 connection to the Loggregator API.
//
// Deprecated: NewV2Client will be removed in the next major version.
// Instead, use v2.NewClient.
func NewV2Client(config Config) (Client, error) {
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

	if config.APIPort != 0 {
		opts = append(opts, v2.WithPort(config.APIPort))
	}

	c, err := v2.NewClient(tlsConfig, opts...)
	if err != nil {
		return nil, err
	}

	return v2shim.NewClient(c), nil
}
