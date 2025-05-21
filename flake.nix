{
  description = "Cablebox control application";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    darwin = {
      url = "github:lnl7/nix-darwin";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs = { self, nixpkgs, flake-utils, darwin }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
        package = pkgs.callPackage ./package.nix { };
      in
      {
        packages = {
          default = package;
          cablebox-control = package;
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
            gotools
          ];
        };
      }) // {
        nixosModule = { config, lib, pkgs, ... }: import ./module.nix { inherit config lib pkgs; package = self.packages.${pkgs.system}.cablebox-control; };
        darwinModule = { config, lib, pkgs, ... }: import ./module.nix { inherit config lib pkgs; package = self.packages.${pkgs.system}.cablebox-control; };
      };
} 