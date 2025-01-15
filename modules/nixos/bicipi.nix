args@{ perSystem, config, flake, lib, pkgs, ... }:
with lib;
let
  cfg = config.services.bicipi;
  bicipi = import ../../packages/bicipi.nix args;
in {
  options = {
    services.bicipi = {
      enable = mkEnableOption "bicipi";
      extraArgs = mkOption {
        description = "Arguments to pass to the bicipi binary";
        default = "";
        example = "--serial=/dev/ttyUSB0 --calibrate=false";
      };
    };
  };
  config = mkIf cfg.enable {
    systemd.services.bicipi = {
      wantedBy = [ "multi-user.target" ];
      after = [ "multi-user.target" ];
      serviceConfig = {
          Restart = "always";
          RestartSec = "10s";
      };
      script = ''
        exec ${bicipi}/bin/bicipi ${cfg.extraArgs}
      '';
    };
  };
}