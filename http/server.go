package httpserver

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v2"
)

var (
	ConfigFile  = flag.String("server.config", "", "path to server configuration file")
	tlsEnable   = flag.Bool("tls", false, "Whether to enable TLS for incoming HTTP requests at configured endpoints. -tlsCertFile and -tlsKeyFile must be set if -tls is set. ")
	tlsCertFile = flag.String("tlsCertFile", "", "Path to file with TLS certificate if -tls is set.")
	tlsKeyFile  = flag.String("tlsKeyFile", "", "Path to file with TLS key if -tls is set.")
)

// Config represents the YAML configuration structure
type Config struct {
	Endpoints  map[string][]EndpointConfig `yaml:"endpoints"`
	TlsConfig  TlsConfig                   `yaml:"tls"`
	FileServer FileServer                  `yaml:"fileServer"`
}

type FileServer struct {
	Enable     bool     `yaml:"enable"`
	ServeDir   string   `yaml:"serveDirectory"`
	ServeFiles []string `yaml:"serveFiles"`
}
type TlsConfig struct {
	ServerCert string `yaml:"serverCert"`
	ServerKey  string `yaml:"serverKey"`
	CaCert     string `yaml:"caCert"`
	EnableMtls bool   `yaml:"enableMtls"`
}

// EndpointConfig represents the configuration for each endpoint
type EndpointConfig struct {
	Method  string `yaml:"method"`
	Content string `yaml:"content"`
	Body    string `yaml:"body"`
	Status  int    `yaml:"status"`
}

// MockServer is a mock HTTP server based on the YAML configuration
type MockServer struct {
	config Config
}

// NewMockServer creates a new instance of MockServer with the provided YAML configuration
func NewMockServer() (*MockServer, error) {
	var cfg Config

	// Read YAML configuration file
	data, err := os.ReadFile(*ConfigFile)
	if err != nil {
		return nil, err
	}
	// Unmarshal YAML data into Config struct
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &MockServer{
		config: cfg,
	}, nil
}

// MetricsHandler handles requests for Prometheus metrics
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}

// Handler handles incoming HTTP requests and returns mock responses based on the configuration
func (s *MockServer) Handler(w http.ResponseWriter, r *http.Request) {
	if !CheckBasicAuth(w, r) {
		return
	}
	endpoint := r.URL.Path
	method := r.Method

	if endpointConfigs, ok := s.config.Endpoints[endpoint]; ok {
		var config *EndpointConfig
		for _, cfg := range endpointConfigs {
			if method == cfg.Method {
				config = &cfg
			}
		}
		if config != nil {
			w.Header().Set("Content-Type", config.Content)
			w.WriteHeader(config.Status)
			fmt.Fprintf(w, config.Body)
			return
		}
	}
	fmt.Fprintf(w, "endpoint not configured\n")
	// Return 404 if the endpoint or method is not configured
	http.NotFound(w, r)
}

func GetListener(addr string) (*net.Listener, error) {
	var tlsConfig *tls.Config
	var cert tls.Certificate
	var err error
	if *tlsEnable {
		cert, err = tls.LoadX509KeyPair(*tlsCertFile, *tlsKeyFile)
		if err != nil {
			log.Printf("unable to load the tls certificates. reason: %s", err.Error())
			return nil, err
		}
		tlsConfig = &tls.Config{
			GetCertificate: func(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
				return &cert, nil
			},
		}
	}
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	if tlsConfig != nil {
		ln = tls.NewListener(ln, tlsConfig)
	}
	return &ln, err
}
