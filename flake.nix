{
  description = "Kin Core Service - Development Environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = {
    self,
    nixpkgs,
    flake-utils,
  }:
    flake-utils.lib.eachDefaultSystem (
      system: let
        pkgs = import nixpkgs {
          inherit system;
        };
      in {
        devShells.default = pkgs.mkShell {
          name = "kin-core-svc";

          buildInputs = with pkgs; [
            # Go
            go

            # Task runner
            go-task

            # Protobuf
            buf
            protobuf
            protoc-gen-go
            protoc-gen-go-grpc

            # Linting & Formatting
            golangci-lint
            gotools # includes goimports

            # Database migrations
            pgroll

            # Development tools
            air # Hot reload
            grpcurl # gRPC testing
            docker-compose

            # Infrastructure tools
            opentofu
            terragrunt
            awscli2
            kubectl
            kubernetes-helm
            argocd
            k9s # Kubernetes TUI

            # Utilities
            git
            jq
          ];

          shellHook = ''
            export GOBIN="$PWD/.nix-go/bin"
            export PATH="$GOBIN:$PATH"

            task setup --silent

            echo "Kin Core Service Development Environment"
            echo ""
            echo "=== Application Tools ==="
            echo "Go:         $(go version | cut -d' ' -f3)"
            echo "Buf:        $(buf --version)"
            echo "Task:       $(task --version)"
            echo ""
            echo "=== Infrastructure Tools ==="
            echo "OpenTofu:   $(tofu version -json | jq -r '.terraform_version')"
            echo "Terragrunt: $(terragrunt --version | head -1 | cut -d' ' -f3)"
            echo "AWS CLI:    $(aws --version | cut -d' ' -f1 | cut -d'/' -f2)"
            echo "kubectl:    $(kubectl version --client -o json 2>/dev/null | jq -r '.clientVersion.gitVersion')"
            echo "Helm:       $(helm version --short | cut -d'+' -f1)"
            echo "ArgoCD:     $(argocd version --client --short 2>/dev/null || echo 'v2.x')"
            echo ""
            echo "Run 'task' to see available commands"
          '';

          # Environment variables
          CGO_ENABLED = "0";
        };
      }
    );
}
