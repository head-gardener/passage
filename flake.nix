{
  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
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
        ...
      }: {
        treefmt = {
          projectRootFile = ".git/config";
          programs.gofmt.enable = true;
          programs.alejandra.enable = true;
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

        packages = {
          passage = with pkgs;
            buildGoModule {
              pname = "passage";
              version = "v0.0.0";

              buildInputs = [self'.packages.bee2];

              src = ./.;

              vendorHash = "sha256-hqcSZ2Peqo7cjQ6+7Ubbhlt8u8autuoVB7BziK0GkKg=";
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
