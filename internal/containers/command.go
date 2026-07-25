package containers

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/briheet/nozarashi/internal/specs"
)

// Base One liners
const (
	ContainerCliName = "container"
)

// Args for system start
var ContainerSystemStartArgs = []string{
	"system",
	"start",
}

// Args for system status
var ContainerSystemStatusArgs = []string{
	"system",
	"status",
}

// Args for system stop
var ContainerSystemStopArgs = []string{
	"system",
	"stop",
}

// ContainerImagePullArgs builds arguments for pulling an OCI image.
func ContainerImagePullArgs(reference string, platform string) []string {
	return []string{
		"image",
		"pull",
		"--platform",
		platform,
		reference,
	}
}

// ContainerImageBuildArgs builds arguments for building a Dockerfile or Containerfile image.
func ContainerImageBuildArgs(name string, filePath string, contextDirectory string) []string {
	return []string{
		"build",
		"--tag",
		name,
		"--file",
		filePath,
		contextDirectory,
	}
}

// ContainerImageLoadArgs builds arguments for loading an OCI image archive.
func ContainerImageLoadArgs(archivePath string) []string {
	return []string{
		"image",
		"load",
		"--input",
		archivePath,
	}
}

// ContainerImageInspectArgs builds arguments for inspecting an image.
func ContainerImageInspectArgs(reference string) []string {
	return []string{
		"image",
		"inspect",
		reference,
	}
}

// ContainerImageDeleteArgs builds arguments for deleting an image.
func ContainerImageDeleteArgs(reference string) []string {
	return []string{
		"image",
		"delete",
		reference,
	}
}

// ContainerVolumeCreateArgs builds arguments for creating a volume.
func ContainerVolumeCreateArgs(name string) []string {
	return []string{
		"volume",
		"create",
		name,
	}
}

// ContainerVolumeInspectArgs builds arguments for inspecting a volume.
func ContainerVolumeInspectArgs(name string) []string {
	return []string{
		"volume",
		"inspect",
		name,
	}
}

// ContainerVolumeDeleteArgs builds arguments for deleting a volume.
func ContainerVolumeDeleteArgs(name string) []string {
	return []string{
		"volume",
		"delete",
		name,
	}
}

// ContainerNetworkCreateArgs builds arguments for creating a network.
func ContainerNetworkCreateArgs(name string) []string {
	return []string{
		"network",
		"create",
		name,
	}
}

// ContainerNetworkInspectArgs builds arguments for inspecting a network.
func ContainerNetworkInspectArgs(name string) []string {
	return []string{
		"network",
		"inspect",
		name,
	}
}

// ContainerNetworkDeleteArgs builds arguments for deleting a network.
func ContainerNetworkDeleteArgs(name string) []string {
	return []string{
		"network",
		"delete",
		name,
	}
}

// ContainerRunArgs builds arguments for creating and starting a service container.
func ContainerRunArgs(
	projectName string,
	containerName string,
	imageName string,
	platform string,
	service *specs.ServiceSpecs,
) []string {
	args := []string{
		"run",
		"--detach",
		"--name",
		containerName,
		"--platform",
		platform,
	}

	// Keep environment arguments stable across runs.
	environmentKeys := make([]string, 0, len(service.Environment))
	for key := range service.Environment {
		environmentKeys = append(environmentKeys, key)
	}
	slices.Sort(environmentKeys)

	for _, key := range environmentKeys {
		args = append(args, "--env", fmt.Sprintf("%s=%s", key, service.Environment[key]))
	}

	// Mount project-scoped named volumes.
	for _, volume := range service.Volumes {
		args = append(
			args,
			"--volume",
			fmt.Sprintf("%s-%s:%s", projectName, volume.Name, volume.Path),
		)
	}

	// Attach project-scoped networks.
	for _, network := range service.Networks {
		args = append(args, "--network", fmt.Sprintf("%s-%s", projectName, network))
	}

	// Publish configured host ports.
	for _, port := range service.Ports {
		args = append(
			args,
			"--publish",
			fmt.Sprintf("%d:%d", port.Host, port.Container),
		)
	}

	// Apple Container accepts one entrypoint executable.
	if len(service.EntryPoint) > 0 {
		args = append(args, "--entrypoint", service.EntryPoint[0])
	}

	args = append(args, imageName)

	// Remaining entrypoint values become arguments to the executable.
	if len(service.EntryPoint) > 1 {
		args = append(args, service.EntryPoint[1:]...)
	}

	args = append(args, service.Command...)
	return args
}

// ContainerInspectArgs builds arguments for inspecting containers.
func ContainerInspectArgs(names ...string) []string {
	args := []string{
		"inspect",
	}

	return append(args, names...)
}

// ContainerStartArgs builds arguments for starting a container.
func ContainerStartArgs(name string) []string {
	return []string{
		"start",
		name,
	}
}

// ContainerStopArgs builds arguments for stopping a container.
func ContainerStopArgs(name string) []string {
	return []string{
		"stop",
		name,
	}
}

// ContainerDeleteArgs builds arguments for deleting a container.
func ContainerDeleteArgs(name string) []string {
	return []string{
		"delete",
		name,
	}
}

// ContainerLogsArgs builds arguments for getting container logs.
func ContainerLogsArgs(name string, number int) []string {
	args := []string{
		"logs",
	}

	if number > 0 {
		args = append(args, "-n", strconv.Itoa(number))
	}

	return append(args, name)
}

// ContainerListArgs builds arguments for listing containers as JSON.
func ContainerListArgs(all bool) []string {
	args := []string{
		"list",
	}

	if all {
		args = append(args, "--all")
	}

	return append(args, "--format", "json")
}

// ContainerExecArgs builds arguments for executing a command in a container.
func ContainerExecArgs(name string, opts ContainerOptions) []string {
	args := []string{
		"exec",
	}

	if opts.Interactive {
		args = append(args, "--interactive")
	}

	if opts.TTY {
		args = append(args, "--tty")
	}

	args = append(args, name)
	return append(args, opts.Args[1:]...)
}
