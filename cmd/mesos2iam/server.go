package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	http_pkg "github.com/schibsted/mesos2iam/http"
	"github.com/schibsted/mesos2iam/pkg"
	"net/http"
	"time"
)

var (
	DEFAULT_LISTENING_IP                 = "0.0.0.0"
	DEFAULT_SERVER_PORT                  = "51679"
	DEFAULT_AWS_CONTAINER_CREDENTIALS_IP = "169.254.170.2"
	// A custom credentials repository for IAM roles
	DEFAULT_CREDENTIALS_URL    = "http://127.0.0.1:8080"
	DEFAULT_MESOS_2_IAM_PREFIX = "TARDIS_SCHID="
)

type Server struct {
	ListeningIp               string
	HostIp                    string
	AppPort                   string
	Verbose                   bool
	AddIPTablesRule           bool
	AwsContainerCredentialsIp string
	CredentialsURL            string
	Mesos2IamPrefix           string
}

func (s *Server) BuildSecurityRequestHandler(dockerClient *docker.Client, credentialsURL string) *http_pkg.SecurityRequestHandler {
	containerRepository := pkg.NewContainerRepository(dockerClient, s.Mesos2IamPrefix)
	pidFinder := pkg.NewPidFinder()

	jobFinder := pkg.NewJobFinder(containerRepository, pidFinder, s.HostIp, s.Mesos2IamPrefix)

	netClient := &http.Client{
		Timeout: time.Second * 10,
	}
	return http_pkg.NewSecurityRequestHandler(jobFinder, netClient, credentialsURL, s.Mesos2IamPrefix)
}

func (s *Server) Run(dockerClient *docker.Client) {
	credentialsRequestHandler := s.BuildSecurityRequestHandler(dockerClient, s.CredentialsURL)
	http.Handle("/v2/credentials", http_pkg.LogHandler(credentialsRequestHandler))

	serverAddr := s.ListeningIp + ":" + s.AppPort
	log.Info("Listening on ", serverAddr)
	log.Info("Host IP: ", s.HostIp)
	log.Panic(http.ListenAndServe(serverAddr, nil))

}

// NewServer will create a new Server with default values.
func NewServer() *Server {
	return &Server{
		DEFAULT_LISTENING_IP,
		"",
		DEFAULT_SERVER_PORT,
		false,
		false,
		DEFAULT_AWS_CONTAINER_CREDENTIALS_IP,
		DEFAULT_CREDENTIALS_URL,
		DEFAULT_MESOS_2_IAM_PREFIX,
	}
}
