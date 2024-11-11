{
  inputs = {
    flake-parts.url = "github:hercules-ci/flake-parts";
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };

  outputs = inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      systems = [ "x86_64-linux" ];
      perSystem = { pkgs, ... }: {
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [ go gopls gotools go-tools ];
        };

        packages.bee2 = with pkgs; stdenv.mkDerivation {
          pname = "bee2";
          version = "v2.1.4";

          src = fetchFromGitHub {
            owner = "agievich";
            repo = "bee2";
            rev = "e0ea53134ff0939857de3e3ead72fa2a5318c6a4";
            hash = "sha256-Fp2mCUknMD7z4WUkzL0E1FaNtW6WVDU2X94+DegRSFc=";
          };

          nativeBuildInputs = [ cmake ];

          doCheck = true;
        };
      };
    };
}
