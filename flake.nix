{
  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [ "x86_64-linux" ];
      perSystem = { pkgs, self', ... }: {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            self'.packages.bee2
            just
          ];
        };

        packages.passage = with pkgs;
          buildGoModule {
            pname = "passage";
            version = "v0.0.0";

            buildInputs = [
              self'.packages.bee2
            ];

            src = ./.;

            vendorHash = "sha256-hqcSZ2Peqo7cjQ6+7Ubbhlt8u8autuoVB7BziK0GkKg=";
          };

        packages.bee2 = with pkgs;
          stdenv.mkDerivation {
            pname = "bee2";
            version = "v2.1.4";

            src = fetchFromGitHub {
              owner = "agievich";
              repo = "bee2";
              rev = "eba3d815b423c9d34a322061c2bec7a09f33d990";
              hash = "sha256-3qkv1ufMORNFdYYoABB+q/d4rxzNlQHDxvtq1rrbReY=";
            };

            nativeBuildInputs = [ cmake ];

            doCheck = true;
          };
      };
    };
}
