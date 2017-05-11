package loggregator

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"code.cloudfoundry.org/go-loggregator/internal/loggregator_v2"
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

type MetronConfig struct {
	UseV2API           bool   `json:"loggregator_use_v2_api"`
	APIPort            int    `json:"loggregator_api_port"`
	CACertPath         string `json:"loggregator_ca_path"`
	CertPath           string `json:"loggregator_cert_path"`
	KeyPath            string `json:"loggregator_key_path"`
	JobDeployment      string `json:"loggregator_job_deployment"`
	JobName            string `json:"loggregator_job_name"`
	JobIndex           string `json:"loggregator_job_index"`
	JobIP              string `json:"loggregator_job_ip"`
	JobOrigin          string `json:"loggregator_job_origin"`
	BatchMaxSize       uint
	BatchFlushInterval time.Duration
}

// NewClient creates a connection to the Loggregator API.
// Users can opt-in to using the v2 API through configuration.
// If an opt-in feature is not required, the v1 and v2 clients are both
// available for public use as well.
func NewClient(config MetronConfig) (Client, error) {
	if config.UseV2API {
		tlsConfig, err := v2.NewTLSConfig(
			config.CACertPath,
			config.CertPath,
			config.KeyPath,
		)
		if err != nil {
			return nil, err
		}

		conn, err := grpc.Dial(
			fmt.Sprintf("localhost:%d", config.APIPort),
			grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		)
		if err != nil {
			return nil, err
		}
		ingressClient := loggregator_v2.NewIngressClient(conn)

		jobOpts := v2.JobOpts{
			Deployment: config.JobDeployment,
			Name:       config.JobName,
			Index:      config.JobIndex,
			IP:         config.JobIP,
			Origin:     config.JobOrigin,
		}

		opts := []v2.V2Option{v2.WithJobOpts(jobOpts)}

		if config.BatchMaxSize != 0 {
			opts = append(opts, v2.WithBatchMaxSize(config.BatchMaxSize))
		}

		if config.BatchFlushInterval != time.Duration(0) {
			opts = append(opts, v2.WithBatchFlushInterval(config.BatchFlushInterval))
		}

		return v2.NewClient(ingressClient, opts...)
	}

	return v1.NewClient()
}
