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
  imports = [
    (mkIf pkgs.stdenv.isDarwin ./darwin.nix)
    (mkIf (!pkgs.stdenv.isDarwin) ./linux.nix)
  ];

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

  config = mkIf cfg.enable {
    # Set default package if not specified
    services.cablebox-control.package = mkDefault pkgs.cablebox-control;
  };
}
