package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.schibsted.io/spt-infrastructure/mesos2iam.git/iptables"
	"os"
)

func main() {
	server := NewServer()
	parseFlags(server)
	setLogLevel(server.Verbose)

	if server.AddIPTablesRule {
		if err := iptables.AddRules(server.AppPort, server.AwsContainerCredentialsIp, server.HostIP); err != nil {
			log.Fatal(err)
		}
	}

	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		log.Panic(err)
		return
	}

	server.Run(dockerClient)
}

func parseFlags(server *Server) {
	flag.BoolVar(&server.Verbose, "verbose", false, "Enable verbosity")
	flag.BoolVar(&server.AddIPTablesRule, "iptables", false, "Add iptables rule (also requires --host-ip)")
	flag.StringVar(&server.HostIP, "host-ip", getFromEnvOrDefault("MESOS2IAM_HOST_IP", DEFAULT_HOST_IP),
		"IP address of host")
	flag.StringVar(&server.AppPort, "app-port",
		getFromEnvOrDefault("MESOS2IAM_SERVER_PORT", DEFAULT_SERVER_PORT),
		"App port")

	flag.StringVar(&server.AwsContainerCredentialsIp, "aws-container-credentials-ip",
		getFromEnvOrDefault("MESOS2IAM_AWS_CONTAINER_CREDENTIALS_IP", DEFAULT_AWS_CONTAINER_CREDENTIALS_IP),
		"IP address of aws container credentials host")
	flag.StringVar(&server.SmaugURL,
		"smaug-url",
		getFromEnvOrDefault("MESOS2IAM_SMAUG_URL", DEFAULT_SMAUG_URL),
		"Smaug Url")
	flag.Parse()
}

func setLogLevel(verbose bool) {
	if verbose {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}
}

func getFromEnvOrDefault(variableName string, defaultValue string) string {
	value := os.Getenv(variableName)
	if value == "" {
		value = defaultValue
	}

	return value
}
