package metrics_test

import (
	"bytes"
	"code.cloudfoundry.org/go-loggregator/metrics"
	"code.cloudfoundry.org/tlsconfig/certtest"
	"crypto/tls"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/prometheus/common/expfmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gogo/protobuf/proto"
	dto "github.com/prometheus/client_model/go"
)

var _ = Describe("PrometheusMetrics", func() {
	var (
		l = log.New(GinkgoWriter, "", log.LstdFlags)
	)

	BeforeEach(func() {
		// This is needed because the prom registry will register
		// the /metrics route with the default http mux which is
		// global
		http.DefaultServeMux = new(http.ServeMux)
	})

	It("serves metrics on a prometheus endpoint", func() {
		r := metrics.NewRegistry(l, metrics.WithServer(0))

		c := r.NewCounter(
			"test_counter",
			metrics.WithMetricTags(map[string]string{"foo": "bar"}),
			metrics.WithHelpText("a counter help text for test_counter"),
		)

		g := r.NewGauge(
			"test_gauge",
			metrics.WithHelpText("a gauge help text for test_gauge"),
			metrics.WithMetricTags(map[string]string{"bar": "baz"}),
		)

		c.Add(10)
		g.Set(10)
		g.Add(1)

		Eventually(func() string { return getMetrics(r.Port()) }).Should(ContainSubstring(`test_counter{foo="bar"} 10`))
		Eventually(func() string { return getMetrics(r.Port()) }).Should(ContainSubstring("a counter help text for test_counter"))
		Eventually(func() string { return getMetrics(r.Port()) }).Should(ContainSubstring(`test_gauge{bar="baz"} 11`))
		Eventually(func() string { return getMetrics(r.Port()) }).Should(ContainSubstring("a gauge help text for test_gauge"))
	})

	It("accepts custom default tags", func() {
		ct := map[string]string{
			"tag": "custom",
		}

		r := metrics.NewRegistry(l, metrics.WithDefaultTags(ct), metrics.WithServer(0))

		r.NewCounter(
			"test_counter",
			metrics.WithHelpText("a counter help text for test_counter"),
		)

		r.NewGauge(
			"test_gauge",
			metrics.WithHelpText("a gauge help text for test_gauge"),
		)

		Eventually(func() string { return getMetrics(r.Port()) }).Should(And(
			ContainSubstring("test_counter"),
			ContainSubstring("test_gauge"),
		))

		metrics := getMetrics(r.Port())
		metricFamilies, err := new(expfmt.TextParser).TextToMetricFamilies(bytes.NewReader([]byte(metrics)))
		Expect(err).ToNot(HaveOccurred())

		for _, family := range metricFamilies {
			for _, metric := range family.GetMetric() {
				Expect(metric.Label).To(ContainElement(
					&dto.LabelPair{Name: proto.String("tag"), Value: proto.String("custom")},
				), fmt.Sprintf("family %s contained a metric without default tags", family.GetName()))
			}
		}
	})

	It("returns the metric when duplicate is created", func() {
		r := metrics.NewRegistry(l, metrics.WithServer(0))

		c := r.NewCounter("test_counter")
		c2 := r.NewCounter("test_counter")

		c.Add(1)
		c2.Add(2)

		Eventually(func() string {
			return getMetrics(r.Port())
		}).Should(ContainSubstring(`test_counter 3`))

		g := r.NewGauge("test_gauge")
		g2 := r.NewGauge("test_gauge")

		g.Add(1)
		g2.Add(2)

		Eventually(func() string {
			return getMetrics(r.Port())
		}).Should(ContainSubstring(`test_gauge 3`))
	})

	It("panics if the metric is invalid", func() {
		r := metrics.NewRegistry(l)

		Expect(func() {
			r.NewCounter("test-counter")
		}).To(Panic())

		Expect(func() {
			r.NewGauge("test-counter")
		}).To(Panic())
	})

	Context("WithTLSServer", func() {
		It("starts a TLS server", func() {
			ca, caFile := generateCA("someCA")
			certFile, keyFile := generateCertKeyPair(ca, "server")

			r := metrics.NewRegistry(
				l,
				metrics.WithTLSServer(0, certFile, keyFile, caFile),
			)

			g := r.NewGauge(
				"test_gauge",
				metrics.WithHelpText("a gauge help text for test_gauge"),
				metrics.WithMetricTags(map[string]string{"bar": "baz"}),
			)
			g.Set(10)

			Eventually(func() string {
				return getMetricsTLS(r.Port(), ca)
			}).Should(ContainSubstring(`test_gauge{bar="baz"} 10`))

			addr := fmt.Sprintf("http://127.0.0.1:%s/metrics", r.Port())
			resp, err := http.Get(addr)
			Expect(err).ToNot(HaveOccurred())
			Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
		})
	})
})

func getMetrics(port string) string {
	addr := fmt.Sprintf("http://127.0.0.1:%s/metrics", port)
	resp, err := http.Get(addr)
	if err != nil {
		return ""
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	return string(respBytes)
}

func getMetricsTLS(port string, ca *certtest.Authority) string {
	caPool, err := ca.CertPool()
	if err != nil {
		log.Fatal(err)
	}

	cert, err := ca.BuildSignedCertificate("client")
	if err != nil {
		log.Fatal(err)
	}

	tlsCert, err := cert.TLSCertificate()
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				Certificates: []tls.Certificate{tlsCert},
				RootCAs:      caPool,
			},
		},
	}

	addr := fmt.Sprintf("https://127.0.0.1:%s/metrics", port)
	resp, err := client.Get(addr)
	if err != nil {
		return ""
	}

	respBytes, err := ioutil.ReadAll(resp.Body)
	Expect(err).ToNot(HaveOccurred())

	return string(respBytes)
}

func generateCA(caName string) (*certtest.Authority, string) {
	ca, err := certtest.BuildCA(caName)
	if err != nil {
		log.Fatal(err)
	}

	caBytes, err := ca.CertificatePEM()
	if err != nil {
		log.Fatal(err)
	}

	fileName := tmpFile(caName+".crt", caBytes)

	return ca, fileName
}

func tmpFile(prefix string, caBytes []byte) string {
	file, err := ioutil.TempFile("", prefix)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = file.Write(caBytes)
	if err != nil {
		log.Fatal(err)
	}

	return file.Name()
}

func generateCertKeyPair(ca *certtest.Authority, commonName string) (string, string) {
	cert, err := ca.BuildSignedCertificate(commonName, certtest.WithDomains(commonName))
	if err != nil {
		log.Fatal(err)
	}

	certBytes, keyBytes, err := cert.CertificatePEMAndPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	certFile := tmpFile(commonName+".crt", certBytes)
	keyFile := tmpFile(commonName+".key", keyBytes)

	return certFile, keyFile
}
