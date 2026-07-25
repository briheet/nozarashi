<!-- # Nozarashi -->

![Nozarashi Ascii](./assets/nozarashi_ascii.png)

Nozarashi orchestrates reproducible [Apple Container](https://github.com/apple/container)
environments from OCI images, Dockerfiles, Containerfiles and Nix inputs.

Define services and resources in `nozarashi.toml`, then manage their images,
containers, networks and volumes as one project.

## Installation

Install using Go:

```sh
go install github.com/briheet/nozarashi
```


## Requirements

- macOS with Apple Container installed
- Nix when building services from Nix inputs
- Go 1.25 or newer when building Nozarashi from source

## Quick start

Build the CLI:

```sh
go build -o bin/nozarashi ./cmd/nozarashi
```

From a directory containing `nozarashi.toml`:

```sh
sudo nozarashi create
nozarashi build
nozarashi up
```

`create` registers the project DNS domain and creates its networks and volumes.
DNS registration requires administrator privileges; see Apple’s
[local DNS documentation](https://github.com/apple/container/blob/main/docs/tutorial.md#set-up-a-local-dns-domain-optional).

Use a service name to target one service:

```sh
nozarashi build backend
nozarashi up backend
nozarashi logs backend
nozarashi exec backend sh
```

## Commands

| Command | Purpose |
| --- | --- |
| `create` | Create project DNS, networks and volumes |
| `build [services...]` | Build or pull service images |
| `up [services...]` | Create and start containers |
| `down [services...]` | Stop containers |
| `restart [services...]` | Restart containers |
| `ps` | List project containers |
| `logs [services...]` | Print container logs |
| `exec <service> <command>` | Run a command in a container |
| `inspect <services...>` | Display container details |
| `destroy` | Delete project containers, images and resources |

See the [configuration specification](docs/SPEC.md) and
[Docker example](examples/simple_docker) for a complete setup.

## License

[MIT](LICENSE)
