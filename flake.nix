{
  description = "A development environment for a Go project";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        devShells.default = pkgs.mkShell {
          packages = [
            pkgs.go
            pkgs.gopls     # lsp
            pkgs.delve     # debugger
          ];

          shellHook = ''
            echo "✅ Go environment loaded."

            # Define absolute paths for Go environment variables
            export PROJECT_ROOT="$(pwd)"
            export GOPATH="''${PROJECT_ROOT}/.go"
            export GOBIN="''${GOPATH}/bin"
            export GOMODCACHE="''${GOPATH}/pkg/mod"
            export GOCACHE="''${GOPATH}/pkg/cache"

            # Create directories if they don't exist
            mkdir -p "''${GOPATH}"
            mkdir -p "''${GOBIN}"
            mkdir -p "''${GOMODCACHE}"
            mkdir -p "''${GOCACHE}"
            export PATH="''${GOBIN}:''${PATH}"

            echo "GOPATH is set to: ''${GOPATH}"
          '';
        };
      }
    );
}
