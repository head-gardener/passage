{
  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nix-filter.url = "github:numtide/nix-filter";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    treefmt-nix.url = "github:numtide/treefmt-nix";

    treefmt-nix.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = inputs @ {flake-parts, ...}:
    flake-parts.lib.mkFlake {inherit inputs;} {
      imports = [inputs.treefmt-nix.flakeModule];
      systems = ["x86_64-linux"];
      perSystem = {
        pkgs,
        self',
        config,
        lib,
        ...
      }: {
        treefmt = {
          projectRootFile = ".git/config";

          programs.gofmt.enable = true;
          programs.alejandra.enable = true;

          settings.formatter.eclint = {
            command = pkgs.eclint;
            options = ["-fix"];
            includes = ["*"];
            excludes = ["*.md"];
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            just
            self'.packages.bee2
            config.treefmt.build.wrapper
          ];
        };

        formatter = config.treefmt.build.wrapper;

        checks = {
          pkg-check = self'.packages.passage.override {
            doCheck = true;
          };
        };

        packages = rec {
          passage = with pkgs;
            lib.makeOverridable
            buildGoModule {
              pname = "passage";
              version = "v0.0.0";

              buildInputs = [bee2];
              subPackages = ["cmd/passage"];

              src = inputs.nix-filter.lib {
                root = inputs.self;
                include = [
                  "go.mod"
                  "go.sum"
                  (inputs.nix-filter.lib.inDirectory ./cmd)
                  (inputs.nix-filter.lib.inDirectory ./pkg)
                  (inputs.nix-filter.lib.inDirectory ./internal)
                ];
              };

              doCheck = false;

              vendorHash = "sha256-lTGn58gUhQrcKiYHhZiUMc9/DwwBrlaCWiBqbeaMpJE=";
            };

          passage-image = pkgs.dockerTools.buildLayeredImage {
            name = "passage";
            tag = "latest";

            contents = with pkgs; [
              coreutils
              dash
              iproute
              iputils
              netcat
              passage
              (pkgs.writeShellScriptBin "passage-wrapped" ''
                #! /usr/bin/env sh
                set -e
                mkdir -p /dev/net && mknod /dev/net/tun c 10 200
                addr="$1"
                shift
                {
                  sleep 1
                  ip a add "$addr/24" dev tun1
                  ip l set dev tun1 up
                } &
                exec passage "$@"
              '')
            ];

            config.Cmd = ["passage"];
          };

          bee2 = with pkgs;
            stdenv.mkDerivation {
              pname = "bee2";
              version = "v2.1.4";

              src = fetchFromGitHub {
                owner = "agievich";
                repo = "bee2";
                rev = "eba3d815b423c9d34a322061c2bec7a09f33d990";
                hash = "sha256-3qkv1ufMORNFdYYoABB+q/d4rxzNlQHDxvtq1rrbReY=";
              };

              nativeBuildInputs = [cmake];

              doCheck = true;
            };
        };
      };
    };
}
