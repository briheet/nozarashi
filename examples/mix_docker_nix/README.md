# Mixed Docker and Nix

Dockerfile backend with PostgreSQL and Redis built from local Nix packages.

Create the DNS domain, volumes and network, then build and start the services.
Apple documents the administrator requirement in
[Set up a local DNS domain](https://github.com/apple/container/blob/main/docs/tutorial.md#set-up-a-local-dns-domain-optional).

```sh
sudo ../../bin/nozarashi create
../../bin/nozarashi build
../../bin/nozarashi up
curl http://localhost:18080/api/v1/health
../../bin/nozarashi destroy
```
