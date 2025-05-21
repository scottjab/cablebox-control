{
  config,
  lib,
  pkgs,
  ...
}:

with lib;

let
  cfg = config.services.cablebox-control;
  package = cfg.package;
in
{
  options.services.cablebox-control = {
    enable = mkEnableOption "Cablebox Control Service";

    package = mkOption {
      type = types.package;
      description = "The cablebox-control package to use.";
    };

    statusSocket = mkOption {
      type = types.path;
      default = "/run/FieldStation42/runtime/play_status.socket";
      description = "Path to the play status socket.";
    };

    channelSocket = mkOption {
      type = types.path;
      default = "/run/FieldStation42/runtime/channel.socket";
      description = "Path to the channel socket.";
    };

    listenAddress = mkOption {
      type = types.str;
      default = "127.0.0.1:8080";
      description = "Address to listen on (e.g. ':8080' or '127.0.0.1:8080').";
    };

    extraArgs = mkOption {
      type = types.listOf types.str;
      default = [ ];
      description = "Extra arguments to pass to cablebox-control.";
    };

    user = mkOption {
      type = types.str;
      default = "cablebox-control";
      description = "User account under which cablebox-control runs.";
    };

    group = mkOption {
      type = types.str;
      default = "cablebox-control";
      description = "Group under which cablebox-control runs.";
    };
  };

  config = mkIf cfg.enable (
    let
      args = [
        "--status-socket=${cfg.statusSocket}"
        "--channel-socket=${cfg.channelSocket}"
        "--listen=${cfg.listenAddress}"
      ] ++ cfg.extraArgs;
    in
    if pkgs.stdenv.isDarwin then
      {
        launchd.user.agents.cablebox-control = {
          enable = true;
          config = {
            Label = "com.github.scottjab.cablebox-control";
            ProgramArguments = [
              "${package}/bin/cablebox-control"
            ] ++ args;
            RunAtLoad = true;
            KeepAlive = true;
            StandardErrorPath = "/tmp/cablebox-control.err.log";
            StandardOutPath = "/tmp/cablebox-control.out.log";
          };
        };
      }
    else
      {
        users.users.${cfg.user} = {
          isSystemUser = true;
          group = cfg.group;
          description = "Cablebox Control Service User";
        };

        users.groups.${cfg.group} = { };

        systemd.services.cablebox-control = {
          description = "Cablebox Control Service";
          after = [ "network.target" ];
          wantedBy = [ "multi-user.target" ];
          serviceConfig = {
            ExecStart = ''
              ${package}/bin/cablebox-control ${concatStringsSep " " args}
            '';
            Restart = "on-failure";
            User = cfg.user;
            Group = cfg.group;
            DynamicUser = false;
            StateDirectory = "cablebox-control";
            StateDirectoryMode = "0750";
          };
        };
      }
  );
}
