package integration_test

import (
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"
	"testing"

	"github.com/briheet/nozarashi/internal/containers"
)

const integrationDomain = "nozarashi-integration"

func TestMain(m *testing.M) {
	containerPath, err := exec.LookPath(containers.ContainerCliName)
	if err != nil {
		os.Exit(m.Run())
	}

	// Authenticate before the lifecycle tests start, while keeping the tests unprivileged.
	sudoValidateCmd := exec.Command("sudo", "-v")
	sudoValidateCmd.Stdin = os.Stdin
	sudoValidateCmd.Stdout = os.Stdout
	sudoValidateCmd.Stderr = os.Stderr
	if err := sudoValidateCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "authenticate integration tests: %v\n", err)
		os.Exit(1)
	}

	// DNS creation is the only integration setup operation that requires root.
	listDNSCmd := exec.Command(containerPath, containers.ContainerSystemDNSListArgs...)
	output, err := listDNSCmd.Output()
	if err != nil {
		fmt.Fprintf(os.Stderr, "list integration DNS domains: %v\n", err)
		os.Exit(1)
	}

	if !slices.Contains(strings.Fields(string(output)), integrationDomain) {
		createDNSCmd := exec.Command(
			"sudo",
			append(
				[]string{"--", containerPath},
				containers.ContainerSystemDNSCreateArgs(integrationDomain)...,
			)...,
		)
		createDNSCmd.Stdin = os.Stdin
		createDNSCmd.Stdout = os.Stdout
		createDNSCmd.Stderr = os.Stderr
		if err := createDNSCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "create integration DNS domain: %v\n", err)
			os.Exit(1)
		}
	}

	os.Exit(m.Run())
}
