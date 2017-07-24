package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/schibsted/mesos2iam/iptables"
	"os"
)

func main() {
	server := NewServer()
	parseFlags(server)

	if server.HostIp == "" {
		log.Panic("HostIp can't be empty")
	}

	setLogLevel(server.Verbose)

	if server.AddIPTablesRule {
		if err := iptables.AddRules(server.AppPort, server.AwsContainerCredentialsIp, server.HostIp); err != nil {
			log.Fatal(err)
		}
	}

	dockerClient, err := docker.NewClientFromEnv()
	if err != nil {
		log.Panic(err)
	}

	server.Run(dockerClient)
}

func parseFlags(server *Server) {
	flag.BoolVar(&server.Verbose, "verbose", false, "Enable verbosity")
	flag.BoolVar(&server.AddIPTablesRule, "iptables", false, "Add iptables rule (also requires --host-ip)")
	flag.StringVar(&server.ListeningIp, "listening-ip", getFromEnvOrDefault("MESOS2IAM_LISTENING_IP", DEFAULT_LISTENING_IP),
		"Listening IP address")
	flag.StringVar(&server.HostIp, "host-ip", getFromEnvOrDefault("MESOS2IAM_HOST_IP", ""),
		"Listening IP address")
	flag.StringVar(&server.AppPort, "app-port",
		getFromEnvOrDefault("MESOS2IAM_SERVER_PORT", DEFAULT_SERVER_PORT),
		"App port")

	flag.StringVar(&server.AwsContainerCredentialsIp, "aws-container-credentials-ip",
		getFromEnvOrDefault("MESOS2IAM_AWS_CONTAINER_CREDENTIALS_IP", DEFAULT_AWS_CONTAINER_CREDENTIALS_IP),
		"IP address of aws container credentials host")
	flag.StringVar(&server.CredentialsURL,
		"credentials-url",
		getFromEnvOrDefault("MESOS2IAM_CREDENTIALS_URL", DEFAULT_CREDENTIALS_URL),
		"Credentials Url")
	flag.StringVar(&server.Mesos2IamPrefix,
		"mesos-2-iam-prefix",
		getFromEnvOrDefault("MESOS2IAM_PREFIX", DEFAULT_MESOS_2_IAM_PREFIX),
		"Mesos2Iam prefix to parse the id to be sent to credentials url")
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
