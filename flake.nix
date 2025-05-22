{
  description = "Cablebox Control Service";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages = {
          cablebox-control = pkgs.buildGoModule {
            pname = "cablebox-control";
            version = "0.1.0";
            src = ./.;
            vendorHash = null; # This will be set automatically on first build
          };
          default = self.packages.${system}.cablebox-control;
        };
      }
    )
    // {
      nixosModules.cablebox-control =
        {
          config,
          lib,
          pkgs,
          ...
        }:
        with lib;
        let
          cfg = config.services.cablebox-control;
        in
        {
          options.services.cablebox-control = {
            enable = mkEnableOption "Cablebox Control Service";
            package = mkOption {
              type = types.package;
              default = self.packages.${pkgs.system}.cablebox-control;
              description = "The cablebox-control package to use.";
            };
            statusSocket = mkOption {
              type = types.str;
              default = "FieldStation42/runtime/play_status.socket";
              description = "Path to the play status socket";
            };
            channelSocket = mkOption {
              type = types.str;
              default = "FieldStation42/runtime/channel.socket";
              description = "Path to the channel socket";
            };
            listenAddress = mkOption {
              type = types.str;
              default = ":8080";
              description = "Address to listen on (e.g. ':8080' or '127.0.0.1:8080')";
            };
          };

          config = mkIf cfg.enable {
            systemd.services.cablebox-control = {
              description = "Cablebox Control Service";
              wantedBy = [ "multi-user.target" ];
              after = [ "network.target" ];
              serviceConfig = {
                ExecStart = "${cfg.package}/bin/cablebox-control -status-socket ${cfg.statusSocket} -channel-socket ${cfg.channelSocket} -listen ${cfg.listenAddress}";
                Restart = "always";
                RestartSec = "10";
                DynamicUser = true;
                ProtectHome = false;
                NoNewPrivileges = false;
              };
            };
          };
        };
    };
}
