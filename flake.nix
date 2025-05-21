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
    let
      supportedSystems = [
        "x86_64-linux"
        "aarch64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];
      perSystem = system:
        let
          pkgs = nixpkgs.legacyPackages.${system};
          package = pkgs.callPackage ./nix/packages { };
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
        };
    in
    flake-utils.lib.eachSystem supportedSystems perSystem
    // {
      nixosModules.default = { config, lib, pkgs, ... }:
        import ./nix/modules/service.nix {
          inherit config lib pkgs;
        };

      darwinModules.default = { config, lib, pkgs, ... }:
        import ./nix/modules/service.nix {
          inherit config lib pkgs;
        };
    };
}
