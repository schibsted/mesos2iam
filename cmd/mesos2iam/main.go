package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	http_pkg "github.schibsted.io/spt-infrastructure/mesos2iam.git/http"
	"github.schibsted.io/spt-infrastructure/mesos2iam.git/pkg"
	"net/http"
	"os"
	"time"
)

var (
	serverAddr             string
	DEFAULT_SERVER_ADDRESS = ":51679"
	verbose                bool
	smaugURL               string
	// Smaug is a credentials repository for IAM roles: https://github.schibsted.io/spt-infrastructure/tardis-smaug
	DEFAULT_SMAUG_URL      = "http://127.0.0.1:8080"
)

func main() {
	parseFlags()
	setLogLevel()

	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		log.Panic(err)
		return
	}

	credentialsRequestHandler := buildSecurityRequestHandler(dockerClient)
	http.Handle("/v2/credentials", http_pkg.LogHandler(credentialsRequestHandler))

	log.Info("Listening on ", serverAddr)
	log.Panic(http.ListenAndServe(serverAddr, nil))
}

func parseFlags() {
	flag.BoolVar(&verbose, "verbose", false, "Enable verbosity")
	flag.StringVar(&serverAddr,
		"server-address",
		getFromEnvOrDefault("MESOS2IAM_SERVER_ADDRESS", DEFAULT_SERVER_ADDRESS),
		"Server address")
	flag.StringVar(&smaugURL,
		"smaug-url",
		getFromEnvOrDefault("MESOS2IAM_SMAUG_URL", DEFAULT_SMAUG_URL),
		"Smaug Url")
	flag.Parse()
}

func setLogLevel() {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func buildSecurityRequestHandler(dockerClient *docker.Client) *http_pkg.SecurityRequestHandler {
	containerRepository := pkg.NewContainerRepository(dockerClient)
	pidFinder := pkg.NewPidFinder()

	jobFinder := pkg.NewJobFinder(containerRepository, pidFinder)

	netClient := &http.Client{
		Timeout: time.Second * 10,
	}
	return http_pkg.NewSecurityRequestHandler(jobFinder, netClient, smaugURL)
}

func getFromEnvOrDefault(variableName string, defaultValue string) string {
	value := os.Getenv(variableName)
	if value == "" {
		value = defaultValue
	}

	return value
}
