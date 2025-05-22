{
  config,
  lib,
  pkgs,
  ...
}:

let
  inherit (lib) mkEnableOption mkOption mkIf;
  inherit (lib.types) str bool path package;
  cfg = config.services.cablebox-control;
  args = [
    "-status-socket" cfg.statusSocket
    "-channel-socket" cfg.channelSocket
    "-listen-address" cfg.listenAddress
  ] ++ cfg.extraArgs;
in
{
  options.services.cablebox-control = {
    enable = mkEnableOption "cablebox-control service";
    package = mkOption {
      type = package;
      default = pkgs.cablebox-control;
      description = "The cablebox-control package to use.";
    };
    statusSocket = mkOption {
      type = str;
      default = "/run/cablebox-control/status.sock";
      description = "Path to the status socket.";
    };
    channelSocket = mkOption {
      type = str;
      default = "/run/cablebox-control/channel.sock";
      description = "Path to the channel socket.";
    };
    listenAddress = mkOption {
      type = str;
      default = "127.0.0.1:8080";
      description = "Address to listen on.";
    };
    extraArgs = mkOption {
      type = lib.types.listOf str;
      default = [ ];
      description = "Extra arguments to pass to the service.";
    };
    user = mkOption {
      type = str;
      default = "cablebox-control";
      description = "User to run the service as.";
    };
    group = mkOption {
      type = str;
      default = "cablebox-control";
      description = "Group to run the service as.";
    };
  };

  config = mkIf cfg.enable (
    if pkgs.stdenv.isDarwin then {
      launchd.user.agents.cablebox-control = {
        enable = true;
        config = {
          Label = "com.github.scottjab.cablebox-control";
          ProgramArguments = [
            "${cfg.package}/bin/cablebox-control"
          ] ++ args;
          RunAtLoad = true;
          KeepAlive = true;
          StandardErrorPath = "/tmp/cablebox-control.err.log";
          StandardOutPath = "/tmp/cablebox-control.out.log";
        };
      };
    } else {
      {
        inherit (config) users systemd;
        users = {
          users.${cfg.user} = {
            isSystemUser = true;
            group = cfg.group;
            description = "Cablebox Control Service User";
          };
          groups.${cfg.group} = { };
        };
        systemd = {
          services.cablebox-control = {
            description = "Cablebox control service";
            wantedBy = [ "multi-user.target" ];
            serviceConfig = {
              ExecStart = "${cfg.package}/bin/cablebox-control ${lib.concatStringsSep " " args}";
              Restart = "always";
              DynamicUser = true;
              RuntimeDirectory = "cablebox-control";
              RuntimeDirectoryMode = "0755";
            };
          };
        };
      }
    }
  );
}
