{
  config,
  lib,
  pkgs,
  ...
}:

with lib;

let
  cfg = config.services.cablebox-control;
  args = [
    "--status-socket=${cfg.statusSocket}"
    "--channel-socket=${cfg.channelSocket}"
    "--listen=${cfg.listenAddress}"
  ] ++ cfg.extraArgs;
in
{
  config = mkIf cfg.enable {
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
  };
} 