package pkg

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/go-errors/errors"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type JobFinder interface {
	FindJobIdFromRequest(request *http.Request) (string, error)
}

func NewJobFinder(repository ContainerRepository, pidFinder PidFinder, hostIp string) JobFinder {
	return &ContainerJobFinder{
		repository,
		pidFinder,
		hostIp,
	}
}

type ContainerJobFinder struct {
	repository ContainerRepository
	pidFinder  PidFinder
	hostIp		string
}

func (finder *ContainerJobFinder) FindJobIdFromRequest(request *http.Request) (jobId string, err error) {
	log.Debugf("Remote address: %s", request.RemoteAddr)
	ip := getIp(request.RemoteAddr)

	containerFinder := finder.buildContainerFinder(ip, request)
	container, err := containerFinder.Find()

	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	jobId, err = DiscoverJobIDFromContainer(container)
	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	return jobId, nil
}

func (finder *ContainerJobFinder) buildContainerFinder(ip string, request *http.Request) (containerFinder ContainerFinder) {
	if ip == finder.hostIp {
		log.Debug("Container in host mode")

		return &ContainerInHostModeFinder{
			finder.repository,
			finder.pidFinder,
			getPort(request.RemoteAddr),
		}
	}
	log.Debug("Container in bridge mode")

	return &ContainerInBridgeModeFinder{
		finder.repository,
		ip,
	}
}

type PidFinder interface {
	GetCommandPidByPort(port string) (int32, error)
}

func NewPidFinder() *CommandPidFinder {
	return &CommandPidFinder{}
}

type CommandPidFinder struct {
}

func (finder *CommandPidFinder) GetCommandPidByPort(port string) (int32, error) {
	command := fmt.Sprintf("sudo fuser %s/tcp", port)

	lsofOutput, err := exec.Command("bash", "-c", command).CombinedOutput()

	if err != nil {
		log.Error(err.Error())
		return 0, err
	}

	pid, err := finder.extractPidFromFuserCommand(string(lsofOutput[:]))

	if err == nil {
		return pid, nil
	}

	return 0, errors.Errorf("Can't get Pid by port %s", port)
}

func (finder *CommandPidFinder) extractPidFromFuserCommand(command_result string) (int32, error) {
	RegexPid := regexp.MustCompile("(?s)^(\\d+)/tcp:\\s+(\\d+)(?s:.*)$")
	match := RegexPid.FindStringSubmatch(command_result)

	if match != nil {
		log.Debug(fmt.Sprintf("Fuser result matches regexp"))
		pid, err := strconv.ParseInt(match[2], 10, 32)

		if err == nil {
			return int32(pid), nil
		}
	}

	return 0, errors.Errorf("Couldn't get PID")
}

func getPort(addr string) string {
	index := strings.Index(addr, ":")

	if index < 0 {
		return addr
	}

	return addr[index+1:]
}

func getIp(addr string) string {
	index := strings.Index(addr, ":")

	if index < 0 {
		return addr
	}

	return addr[:index]
}
