package iptables

import (
	"errors"
	"github.com/coreos/go-iptables/iptables"
)

// AddRules adds the required rules to the host's nat table
func AddRules(appPort, metadataAddress, hostIP string) error {
	if hostIP == "" {
		return errors.New("--host-ip must be set")
	}

	ipt, err := iptables.New()
	if err != nil {
		return err
	}

	// These rules must be the first, otherwise, docker rules will override them.
	if err := ipt.Insert(
		"nat", "PREROUTING", 0, "-p", "tcp", "-d", metadataAddress, "--dport", "80",
		"-j", "DNAT", "--to-destination", hostIP+":"+appPort,
	); err != nil {
		return err
	}

	if err := ipt.Insert(
		"nat", "OUTPUT", 0, "-p", "tcp", "-m", "tcp", "-d", metadataAddress, "--dport", "80",
		"-j", "REDIRECT", "--to-ports", appPort,
	); err != nil {
		return err
	}

	return nil
}
