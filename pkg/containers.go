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
