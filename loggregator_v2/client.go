package loggregator_v2

import (
	"time"

	"code.cloudfoundry.org/lager"

	"github.com/cloudfoundry/sonde-go/events"
)

//go:generate bash scripts/generate_protos.sh

//go:generate counterfeiter -o fakes/fake_client.go . Client

//go:generate counterfeiter -o fakes/fake_batcher.go . Batcher
type Batcher interface {
	ComponentMetricsClient
	Send() error
}

type ComponentClient interface {
	ComponentMetricsClient
	Batcher() Batcher
	IncrementCounter(name string) error
}

type ComponentMetricsClient interface {
	SendDuration(name string, value time.Duration) error
	SendMebiBytes(name string, value int) error
	SendMetric(name string, value int) error
	SendBytesPerSecond(name string, value float64) error
	SendRequestsPerSecond(name string, value float64) error
}

type Client interface {
	ComponentClient
	SendAppLog(appID, message, sourceType, sourceInstance string) error
	SendAppErrorLog(appID, message, sourceType, sourceInstance string) error
	SendAppMetrics(metrics *events.ContainerMetric) error
}

type MetronConfig struct {
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
	DropsondePort int    `json:"dropsonde_port"`
}

func NewClient(logger lager.Logger, config MetronConfig) (Client, error) {
	if config.UseV2API {
		return newGrpcClient(logger, config)
	}

	return newDropsondeClient()
}
