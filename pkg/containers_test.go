package pkg_test

import (
	"github.com/fsouza/go-dockerclient"
	"github.com/stretchr/testify/assert"
	"github.com/schibsted/mesos2iam/pkg"
	"testing"
)

func TestDiscoverApplicationFromContainerReturnsApplicationName(t *testing.T) {
	container := &docker.Container{
		Config: &docker.Config{
			Env: []string{"TARDIS_SCHID=4ea13548-caa8-48dc-af69-58a651d9fa3b"},
		},
	}

	applicationName, err := pkg.DiscoverJobIDFromContainer(container, "TARDIS_SCHID=")

	assert.Equal(t, err, nil)
	assert.Equal(t, "4ea13548-caa8-48dc-af69-58a651d9fa3b", applicationName)
}

func TestDiscoverApplicationFromContainerReturnsErrorIfMesosTaskIdDoesNotExist(t *testing.T) {
	container := &docker.Container{
		Config: &docker.Config{
			Env: []string{},
		},
	}

	applicationName, err := pkg.DiscoverJobIDFromContainer(container, "TARDIS_SCHID=")

	if assert.Error(t, err, "An error was expected if TARDIS_SCHID envvar does not exist") {
		assert.Equal(t, err.Error(), "Couldn't get TARDIS_SCHID environment variable from container")
	}
	assert.Equal(t, "", applicationName)
}

func TestDiscoverApplicationFromContainerReturnsErrorIfCanNotGetApplicationName(t *testing.T) {
	container := &docker.Container{
		Config: &docker.Config{
			Env: []string{"TARDIS_SCHID=stupidcontent"},
		},
	}

	applicationName, err := pkg.DiscoverJobIDFromContainer(container, "TARDIS_SCHID=")

	if assert.Error(t, err, "An error was expected if TARDIS_SCHID does not contain a valid uuid") {
		assert.Equal(t, err.Error(), "SCHID \"stupidcontent\" is not a valid uuidv4")
	}
	assert.Equal(t, "", applicationName)
}
