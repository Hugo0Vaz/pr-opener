{
  description = "Go development environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          config.allowUnfree = true;
        };

        goVersion = pkgs.go_1_24;

        devTools = with pkgs; [
          goVersion
          gopls            # Official Go language server
          delve            # Debugger
          go-tools         # Static analysis (staticcheck, etc)
          golangci-lint    # Popular linter aggregator
          gomodifytags     # Go struct tags manipulation
          gotests          # Generate tests
          richgo           # Colorized `go test` output
          mockgen          # Generate mocks
          air              # Live reload for Go apps
          gofumpt          # Strict gofmt
          revive           # Fast linter
          govulncheck      # Vulnerability scanner
          sqlc             # SQL to Go type-safe codegen
          protobuf         # Protocol Buffers
          grpcurl          # gRPC debugging
        ];
      in
      {
        devShells.default = pkgs.mkShell {
          name = "go-dev-env";

          packages = devTools;

          shellHook = ''
            echo "Go `${pkgs.go}/bin/go version`"
            echo "Available tools:"
            echo "  gopls, delve, golangci-lint, air, sqlc, etc."
          '';
        };
      });
}
