package iptables

import (
	"errors"
	"github.com/coreos/go-iptables/iptables"
)

// AddRules adds the required rules to the host's nat table
func AddRules(appPort, metadataAddress, hostIp string) error {
	if hostIp == "" {
		return errors.New("--host-ip must be set")
	}

	ipt, err := iptables.New()
	if err != nil {
		return err
	}

	table := "nat"

	rulespec := []string{"-p", "tcp",
			     "-d", metadataAddress,
			     "--dport", "80",
			     "-j", "DNAT",
			     "--to-destination", hostIp + ":" + appPort}
	err = insertIfDoesNotExist(ipt, table, "PREROUTING", 1, rulespec)
	if err != nil {
		return err
	}

	rulespec = []string{"-p", "tcp",
			    "-m", "tcp",
			    "-d", metadataAddress,
			    "--dport", "80",
			    "-j", "REDIRECT", "--to-ports", appPort}
	err = insertIfDoesNotExist(ipt, table, "OUTPUT", 1, rulespec)
	if err != nil {
		return err
	}

	return nil
}

func insertIfDoesNotExist(ipt *iptables.IPTables, table string, chain string, pos int, rulespec []string) error {
	// These rules must be the first, otherwise, docker rules will override them.
	if exists, _ := ipt.Exists(table, chain, rulespec...); !exists {
		if err := ipt.Insert(table, chain, pos, rulespec...); err != nil {
			return err
		}
	}

	return nil
}
