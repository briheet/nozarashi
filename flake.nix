{
  description = "Orchestrate Apple Container environments from OCI, Dockerfile and Nix inputs";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

    systems = {
      url = "github:nix-systems/default";
      flake = false;
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      systems,
    }:
    let
      inherit (nixpkgs) lib legacyPackages;

      forAllSystems = lib.genAttrs (import systems);
    in
    {
      packages = forAllSystems (
        system:
        let
          pkgs = legacyPackages.${system};
        in
        rec {
          nozarashi = pkgs.buildGoModule {
            pname = "nozarashi";
            version = "0.1.0";

            src = self;
            vendorHash = "sha256-HOH/N4c2CKo+2CKypnuJIe5s/JGebmWlmaFqXPjGlqg=";

            subPackages = [ "cmd/nozarashi" ];
            ldflags = [
              "-s"
              "-w"
            ];

            meta = {
              description = "Orchestrate Apple Container development environments";
              homepage = "https://github.com/briheet/nozarashi";
              license = lib.licenses.mit;
              mainProgram = "nozarashi";
              platforms = lib.platforms.darwin;
            };
          };

          default = nozarashi;
        }
      );

      apps = forAllSystems (system: {
        default = {
          type = "app";
          program = lib.getExe self.packages.${system}.default;
          meta.description = "Run Nozarashi";
        };
      });

      devShells = forAllSystems (
        system:
        let
          pkgs = legacyPackages.${system};
        in
        {
          default = pkgs.mkShell {
            packages = [
              pkgs.go
              pkgs.nix
            ]
            ++ lib.optionals pkgs.stdenv.hostPlatform.isDarwin [
              pkgs.container
            ];
          };
        }
      );

      formatter = forAllSystems (system: legacyPackages.${system}.nixfmt);
    };
}
