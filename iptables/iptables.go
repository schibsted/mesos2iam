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

	if err := ipt.AppendUnique(
		"nat", "PREROUTING", "-p", "tcp", "-d", metadataAddress, "--dport", "80",
		"-j", "DNAT", "--to-destination", hostIP+":"+appPort,
	); err != nil {
		return err
	}

	if err := ipt.AppendUnique(
		"nat", "OUTPUT", "-p", "tcp", "-m", "tcp", "-d", metadataAddress, "--dport", "80",
		"-j", "REDIRECT", "--to-ports", appPort,
	); err != nil {
		return err
	}

	return nil
}
