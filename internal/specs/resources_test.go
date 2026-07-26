package specs

import (
	"encoding/json"
	"testing"
)

func TestDecodeRuntimeResources(t *testing.T) {
	var images []Image
	if err := json.Unmarshal([]byte(`[{
		"id":"abcdef1234567890",
		"configuration":{
			"name":"example/api:latest",
			"creationDate":"2026-07-26T08:47:13Z"
		},
		"variants":[{
			"platform":{"architecture":"arm64","os":"linux"},
			"size":25896597
		}]
	}]`), &images); err != nil {
		t.Fatalf("decode images: %v", err)
	}
	if len(images) != 1 ||
		images[0].Configuration.Name != "example/api:latest" ||
		images[0].Configuration.CreationDate != "2026-07-26T08:47:13Z" ||
		images[0].Variants[0].Platform.Architecture != "arm64" ||
		images[0].Variants[0].Platform.OS != "linux" ||
		images[0].Variants[0].Size != 25896597 {
		t.Fatalf("unexpected images: %#v", images)
	}

	var volumes []Volume
	if err := json.Unmarshal([]byte(`[{
		"id":"example-data",
		"configuration":{
			"creationDate":"2026-07-26T08:47:13Z",
			"driver":"local",
			"format":"ext4",
			"sizeInBytes":549755813888
		}
	}]`), &volumes); err != nil {
		t.Fatalf("decode volumes: %v", err)
	}
	if len(volumes) != 1 ||
		volumes[0].Configuration.CreationDate != "2026-07-26T08:47:13Z" ||
		volumes[0].Configuration.Driver != "local" ||
		volumes[0].Configuration.Format != "ext4" ||
		volumes[0].Configuration.SizeInBytes != 549755813888 {
		t.Fatalf("unexpected volumes: %#v", volumes)
	}

	var networks []Network
	if err := json.Unmarshal([]byte(`[{
		"id":"example-default",
		"configuration":{"mode":"nat"},
		"status":{
			"ipv4Gateway":"192.168.68.1",
			"ipv4Subnet":"192.168.68.0/24",
			"ipv6Subnet":"fd80:1f45:6fad:5fe7::/64"
		}
	}]`), &networks); err != nil {
		t.Fatalf("decode networks: %v", err)
	}
	if len(networks) != 1 ||
		networks[0].Configuration.Mode != "nat" ||
		networks[0].Status.IPv4Gateway != "192.168.68.1" ||
		networks[0].Status.IPv4Subnet != "192.168.68.0/24" ||
		networks[0].Status.IPv6Subnet != "fd80:1f45:6fad:5fe7::/64" {
		t.Fatalf("unexpected networks: %#v", networks)
	}
}
