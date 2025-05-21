{ config, lib, pkgs, ... }:

with lib;

let
  cfg = config.services.cablebox-control;
  package = cfg.package;
  isDarwin = pkgs.stdenv.isDarwin;
in {
  options.services.cablebox-control = {
    enable = mkOption {
      type = types.bool;
      default = false;
      description = "Enable the cablebox-control service.";
    };
    package = mkOption {
      type = types.package;
      default = pkgs.cablebox-control or pkgs.callPackage ./default.nix { };
      defaultText = literalExpression "pkgs.cablebox-control";
      description = "The cablebox-control package to use.";
    };
    extraArgs = mkOption {
      type = types.listOf types.str;
      default = [];
      description = "Extra arguments to pass to cablebox-control.";
    };
  };

  config = mkIf cfg.enable (
    if isDarwin then {
      launchd.user.agents.cablebox-control = {
        enable = true;
        config = {
          Label = "com.github.scottjab.cablebox-control";
          ProgramArguments = [
            "${package}/bin/cablebox-control"
          ] ++ cfg.extraArgs;
          RunAtLoad = true;
          KeepAlive = true;
          StandardErrorPath = "/tmp/cablebox-control.err.log";
          StandardOutPath = "/tmp/cablebox-control.out.log";
        };
      };
    } else {
      systemd.services.cablebox-control = {
        description = "Cablebox Control Service";
        after = [ "network.target" ];
        wantedBy = [ "multi-user.target" ];
        serviceConfig = {
          ExecStart = ''
            ${package}/bin/cablebox-control ${concatStringsSep " " cfg.extraArgs}
          '';
          Restart = "on-failure";
          User = "cablebox-control";
          Group = "cablebox-control";
        };
        preStart = ''
          getent passwd cablebox-control > /dev/null || useradd -r -s /sbin/nologin cablebox-control
        '';
      };
    }
  );
} 