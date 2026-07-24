package specs

import "runtime"

// PlatformSpecs defines the Linux platform used for building container images.
type PlatformSpecs struct {
	// Operating system written to the OCI platform.
	OperatingSystem string

	// Nix system used for resolving packages.
	NixSystem string

	// OCI architecture written to the image metadata.
	Architecture string
}

// CurrentPlatform contains the Linux platform matching the host architecture.
var CurrentPlatform PlatformSpecs

func init() {
	// Apple Containers uses Linux arm64 images by default.
	CurrentPlatform = PlatformSpecs{
		OperatingSystem: "linux",
		NixSystem:       "aarch64-linux",
		Architecture:    "arm64",
	}

	// Intel hosts use Linux amd64 images.
	if runtime.GOARCH == "amd64" {
		CurrentPlatform = PlatformSpecs{
			OperatingSystem: "linux",
			NixSystem:       "x86_64-linux",
			Architecture:    "amd64",
		}
	}
}
