{ config, lib, ... }:
with lib;
let
  cfg = config.services.bicipi;
  bicipi = ./../../packages/bicipi.nix;
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
      script = ''
        exec ${bicipi}/bin/bicipi ${cfg.extraArgs}
      '';
    };
  };
}