package specs

// Image contains the image list fields used by the TUI.
type Image struct {
	ID            string             `json:"id"`
	Configuration ImageConfiguration `json:"configuration"`
	Variants      []ImageVariant     `json:"variants"`
}

// ImageConfiguration contains the image reference.
type ImageConfiguration struct {
	Name         string `json:"name"`
	CreationDate string `json:"creationDate"`
}

// ImageVariant contains the size of one image manifest.
type ImageVariant struct {
	Platform ImagePlatform `json:"platform"`
	Size     int64         `json:"size"`
}

// ImagePlatform identifies the operating system and architecture of an image variant.
type ImagePlatform struct {
	Architecture string `json:"architecture"`
	OS           string `json:"os"`
}

// Volume contains a managed volume and its persistent configuration.
type Volume struct {
	ID            string              `json:"id"`
	Configuration VolumeConfiguration `json:"configuration"`
}

// VolumeConfiguration contains the volume list fields used by the TUI.
type VolumeConfiguration struct {
	CreationDate string `json:"creationDate"`
	Driver       string `json:"driver"`
	Format       string `json:"format"`
	SizeInBytes  uint64 `json:"sizeInBytes"`
}

// Network contains the network list fields used by the TUI.
type Network struct {
	ID            string               `json:"id"`
	Configuration NetworkConfiguration `json:"configuration"`
	Status        NetworkStatus        `json:"status"`
}

// NetworkConfiguration contains the network mode.
type NetworkConfiguration struct {
	Mode string `json:"mode"`
}

// NetworkStatus contains the network's assigned addresses.
type NetworkStatus struct {
	IPv4Gateway string `json:"ipv4Gateway"`
	IPv4Subnet  string `json:"ipv4Subnet"`
	IPv6Subnet  string `json:"ipv6Subnet"`
}
