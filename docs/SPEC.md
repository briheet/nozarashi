# Nozarashi Configuration Specification

The Go types in [`internal/specs/specs.go`](../internal/specs/specs.go) are the source of truth for this configuration model. The file is decoded using the `mapstructure` tags below.

## Top-level structure

```toml
[project]
[inputs.<name>]
[services.<name>]
[volumes.<name>]
[networks.<name>]
[configs.<name>]
[secrets.<name>]
```

`project`, `inputs`, `services`, `volumes`, and `networks` are required. `configs` and `secrets` are optional.

## Project

```toml
[project]
name = "nozarashi-example"
version = "0.1.0"
description = "Example development environment"
profile = "dev"
```

## Inputs

Inputs may be a Nix flake, Git repository, or local path.

```toml
[inputs.nixpkgs]
type = "flake"
source = "github:NixOS/nixpkgs/nixos-unstable"
ref = "nixos-unstable"

[inputs.backend]
type = "git"
source = "https://github.com/example/backend.git"
ref = "v1.4.0"

[inputs.local]
type = "path"
source = "./modules/local"
```

## Services

Each service describes a long-running workload. `type` identifies how its source is resolved; `reference` points to that source and `attribute` optionally selects an attribute from it.

```toml
[services.api]
type = "path"
reference = "./api"
attribute = "default"
entrypoint = ["/usr/bin/env"]
command = ["./api", "serve"]
environment = { PORT = "1000" }
depends_on = ["postgres"]
replicas = 1
policy = "on-failure"
networks = ["default"]

[[services.api.volumes]]
name = "api-data"
path = "/var/lib/api"

[[services.api.ports]]
host = 1000
container = 1000
```

Supported service source types currently include `image`, `path`, and `input`.

## Volumes and networks

```toml
[volumes.postgres-data]
driver = "local"

[networks.default]
driver = "default"
```

## Configs and secrets

```toml
[configs.api]
name = "api-config"
file = "./config/api.toml"
external = false

[secrets.database]
name = "database-password"
file = "./secrets/database.txt"
external = false
```

Configs may use either `file` or inline `content`. Secrets support `file`; external values set `external = true`.
