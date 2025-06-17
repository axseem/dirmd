{
  description = "A CLI tool for bundling directories into a single markdown file";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        pname = "dirmd";
        version = "0.1.0";
      in
      {
        packages = {
          dirmd = pkgs.buildGoModule {
            inherit pname version;
            src = self;
            vendorHash = "sha256-lVUSnffsoyE1Rn8Le66CeynR/lID0hDmodE6JWt/Nao=";
            subPackages = [ "." ];
          };

          default = self.packages.${system}.dirmd;
        };

        apps = {
          dirmd = {
            type = "app";
            program = "${self.packages.${system}.dirmd}/bin/dirmd";
          };
          default = self.apps.${system}.dirmd;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.go
            pkgs.gopls
            pkgs.gofumpt
          ];
        };
      }
    );
}