package main

import (
	"github.schibsted.io/spt-infrastructure/mesos2iam.git/pkg"
	"time"
	"github.com/fsouza/go-dockerclient"
	http_pkg "github.schibsted.io/spt-infrastructure/mesos2iam.git/http"
	"net/http"
	log "github.com/Sirupsen/logrus"
)

var(
	DEFAULT_HOST_IP                      = "0.0.0.0"
	DEFAULT_SERVER_PORT                  = "51679"
	DEFAULT_AWS_CONTAINER_CREDENTIALS_IP = "169.254.170.2"
	// Smaug is a credentials repository for IAM roles: https://github.schibsted.io/spt-infrastructure/tardis-smaug
	DEFAULT_SMAUG_URL = "http://127.0.0.1:8080"
)

type Server struct {
	HostIP                    string
	AppPort                   string
	Verbose                   bool
	AddIPTablesRule           bool
	AwsContainerCredentialsIp string
	SmaugURL                  string
}

func (s *Server) BuildSecurityRequestHandler(dockerClient *docker.Client, smaugURL string) *http_pkg.SecurityRequestHandler {
	containerRepository := pkg.NewContainerRepository(dockerClient)
	pidFinder := pkg.NewPidFinder()

	jobFinder := pkg.NewJobFinder(containerRepository, pidFinder)

	netClient := &http.Client{
		Timeout: time.Second * 10,
	}
	return http_pkg.NewSecurityRequestHandler(jobFinder, netClient, smaugURL)
}


func (s *Server) Run(dockerClient *docker.Client) {
	credentialsRequestHandler := s.BuildSecurityRequestHandler(dockerClient, s.SmaugURL)
	http.Handle("/v2/credentials", http_pkg.LogHandler(credentialsRequestHandler))

	serverAddr := s.HostIP + ":" + s.AppPort
	log.Info("Listening on ", serverAddr)
	log.Panic(http.ListenAndServe(serverAddr, nil))

}

// NewServer will create a new Server with default values.
func NewServer() *Server {
	return &Server{
		DEFAULT_HOST_IP,
		DEFAULT_SERVER_PORT,
		false,
		false,
		DEFAULT_AWS_CONTAINER_CREDENTIALS_IP,
		DEFAULT_SMAUG_URL,
	}
}
