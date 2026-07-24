package nix

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/briheet/nozarashi/internal/specs"
)

// BuildServiceImage builds a layered OCI image archive for a Nix-backed service.
func BuildServiceImage(ctx context.Context, serviceName string, input specs.InputSpecs, service *specs.ServiceSpecs) (string, error) {
	// Build the expression that resolves the input package.
	inputExpression, err := buildInputExpression(input, service.Attribute)
	if err != nil {
		return "", fmt.Errorf("build Nix input expression for service %q: %w", serviceName, err)
	}

	// Build environment values in a stable order.
	environmentKeys := make([]string, 0, len(service.Environment))
	for key := range service.Environment {
		environmentKeys = append(environmentKeys, key)
	}
	slices.Sort(environmentKeys)

	environment := make([]string, 0, len(environmentKeys))
	for _, key := range environmentKeys {
		environment = append(environment, fmt.Sprintf("%s=%s", key, service.Environment[key]))
	}

	// Generate the dockerTools expression and convert its archive to OCI.
	expression := fmt.Sprintf(
		BuildServiceImageExpression,
		strconv.Quote(specs.CurrentPlatform.NixSystem),
		strconv.Quote(serviceName),
		inputExpression,
		strconv.Quote(specs.CurrentPlatform.Architecture),
		buildNixList(service.EntryPoint),
		buildNixList(service.Command),
		buildNixList(environment),
	)

	// Execute Nix and return the generated image archive path.
	buildCmd := exec.CommandContext(
		ctx,
		NixCliName,
		NixBuildExpressionArgs(expression)...,
	)
	buildCmd.Stderr = os.Stderr

	output, err := buildCmd.Output()
	if err != nil {
		return "", fmt.Errorf("build Nix image for service %q: %w", serviceName, err)
	}

	return strings.TrimSpace(string(output)), nil
}

// buildInputExpression creates the Nix expression used to resolve a package attribute.
func buildInputExpression(input specs.InputSpecs, attribute string) (string, error) {
	switch input.Type {
	case specs.InputTypeNix:
		return fmt.Sprintf(
			FlakeInputExpression,
			strconv.Quote(input.Source),
			strconv.Quote(attribute),
		), nil

	case specs.InputTypeGit:
		ref := ""
		if input.Ref != "" {
			ref = fmt.Sprintf(GitInputRefExpression, strconv.Quote(input.Ref))
		}

		return fmt.Sprintf(
			GitInputExpression,
			strconv.Quote(input.Source),
			ref,
			strconv.Quote(attribute),
		), nil

	case specs.InputTypeLocal:
		source, err := filepath.Abs(input.Source)
		if err != nil {
			return "", fmt.Errorf("resolve local Nix input %q: %w", input.Source, err)
		}

		return fmt.Sprintf(
			LocalInputExpression,
			strconv.Quote(source),
			strconv.Quote(attribute),
		), nil

	default:
		return "", fmt.Errorf("unsupported Nix input type %q", input.Type)
	}
}

// buildNixList converts string values into a Nix list.
func buildNixList(values []string) string {
	items := make([]string, 0, len(values))
	for _, value := range values {
		items = append(items, strconv.Quote(value))
	}

	return fmt.Sprintf("[ %s ]", strings.Join(items, " "))
}
