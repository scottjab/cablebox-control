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
          ${cfg.package}/bin/cablebox-control ${concatStringsSep " " args}
        '';
        Restart = "on-failure";
        User = cfg.user;
        Group = cfg.group;
        DynamicUser = false;
        StateDirectory = "cablebox-control";
        StateDirectoryMode = "0750";
      };
    };
  };
} 