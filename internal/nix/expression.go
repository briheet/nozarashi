package nix

// FlakeInputExpression resolves a package from a flake input.
const FlakeInputExpression = `
  input = builtins.getFlake %s;

  package = builtins.getAttr %s (
    builtins.getAttr system input.packages
  );
`

// GitInputRefExpression adds an optional Git reference.
const GitInputRefExpression = `
    ref = %s;`

// GitInputExpression resolves a package from a Git input.
const GitInputExpression = `
  source = builtins.fetchGit {
    url = %s;%s
  };

  imported = import source;

  input =
    if builtins.isFunction imported
    then imported { inherit pkgs; }
    else imported;

  package = builtins.getAttr %s input;
`

// LocalInputExpression resolves a package from a local Nix input.
const LocalInputExpression = `
  imported = import (builtins.toPath %s);

  input =
    if builtins.isFunction imported
    then imported { inherit pkgs; }
    else imported;

  package = builtins.getAttr %s input;
`

// BuildServiceImageExpression builds a layered image and converts it to an OCI archive.
const BuildServiceImageExpression = `
let
  system = %s;
  nativePkgs = import <nixpkgs> { };
  pkgs = import <nixpkgs> { inherit system; };
  imageName = %s;

%s

  dockerImage = nativePkgs.dockerTools.buildLayeredImage {
    name = imageName;
    tag = "latest";
    architecture = %s;
    contents = [ package ];

    config = {
      Entrypoint = %s;
      Cmd = %s;
      Env = %s;
    };
  };
in
nativePkgs.runCommand "${imageName}-oci.tar" {
  nativeBuildInputs = [ nativePkgs.skopeo ];
} ''
  skopeo --insecure-policy copy \
    docker-archive:${dockerImage} \
    oci-archive:$out:${imageName}:latest
''
`
