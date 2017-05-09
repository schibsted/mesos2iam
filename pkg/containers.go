package pkg

import (
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/go-errors/errors"
	"github.com/shirou/gopsutil/process"
	"regexp"
	"strings"
)

var (
	TARDIS_SCHID_PREFIX = "TARDIS_SCHID="
)

type ContainerRepository interface {
	FindContainerUsingCommandPID(pid int32) (*docker.Container, error)
	FindContainerUsingIp(ip string) (*docker.Container, error)
}

func NewContainerRepository(client *docker.Client) *DockerContainerRepository {
	return &DockerContainerRepository{
		client,
	}
}

// implements ContainerRepository
type DockerContainerRepository struct {
	docker *docker.Client
}

func (repository *DockerContainerRepository) findByContainerPID(pid int32) (*docker.Container, error) {
	containers, err := repository.docker.ListContainers(docker.ListContainersOptions{})
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	for _, container := range containers {
		containerInfo, err := repository.docker.InspectContainer(container.ID)
		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		if pid == int32(containerInfo.State.Pid) {
			log.Debug("Found PID: ", pid)

			return containerInfo, nil
		}
	}

	return nil, errors.Errorf("Container that contains process %d does not exist", pid)
}

func (repository *DockerContainerRepository) FindContainerUsingCommandPID(pid int32) (*docker.Container, error) {
	proc, err := process.NewProcess(int32(pid))

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	parent, err := proc.Parent()

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	container, err := repository.findByContainerPID(parent.Pid)

	if err != nil {
		log.Error(err)
		return nil, errors.Errorf("Container that contains process %d does not exist", pid)
	}

	return container, nil
}

func (repository *DockerContainerRepository) FindContainerUsingIp(ip string) (*docker.Container, error) {
	containers, err := repository.docker.ListContainers(docker.ListContainersOptions{})

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	for _, container := range containers {
		containerInfo, err := repository.docker.InspectContainer(container.ID)

		if err != nil {
			log.Error(err.Error())
			return nil, err
		}

		if ip == containerInfo.NetworkSettings.IPAddress {
			log.Debug("Found IP: ", ip)

			return containerInfo, nil
		}
	}

	return nil, errors.Errorf("Container with ip %s does not exist", ip)
}

type ContainerFinder interface {
	Find() (*docker.Container, error)
}

type ContainerInHostModeFinder struct {
	repository ContainerRepository
	pidFinder  PidFinder
	port       string
}

func (finder *ContainerInHostModeFinder) Find() (*docker.Container, error) {
	log.Debug("Remote port: ", finder.port)

	pid, err := finder.pidFinder.GetCommandPidByPort(finder.port)
	log.Debug("Pid: ", pid)

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	container, err := finder.repository.FindContainerUsingCommandPID(pid)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return container, err
}

type ContainerInBridgeModeFinder struct {
	repository ContainerRepository
	ip         string
}

func (finder *ContainerInBridgeModeFinder) Find() (*docker.Container, error) {
	container, err := finder.repository.FindContainerUsingIp(finder.ip)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return container, err
}

func findJobID(container *docker.Container) (string, error) {
	for _, envvar := range container.Config.Env {
		if strings.HasPrefix(envvar, TARDIS_SCHID_PREFIX) {
			return strings.TrimPrefix(envvar, TARDIS_SCHID_PREFIX), nil
		}
	}

	return "", errors.Errorf("Couldn't get TARDIS_SCHID environment variable from container")
}

func DiscoverJobIDFromContainer(container *docker.Container) (string, error) {
	jobID, err := findJobID(container)
	if err != nil {
		log.Error(err)
		return "", err
	}

	if isValidUUID(jobID) {
		return jobID, nil
	}

	return "", errors.Errorf("SCHID \"%s\" is not a valid uuidv4", jobID)
}

func isValidUUID(uuid string) bool {
	r := regexp.MustCompile("^[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}$")
	return r.MatchString(uuid)
}
