# resizes a partition to occupy empty space
# boot.growPartition=true will grow the partition at /
# boot.growPartitionCustom allows the device to be specified

{ config, lib, pkgs, ... }:
with lib;
let
  cfg = config.boot.growPartitionCustom;
in {
  options = {
    boot.growPartitionCustom = {
      enable = mkEnableOption "growPartitionCustom";
      device = mkOption {
        description = "The partition device to grow";
        default = "/dev/root";
        example = "/dev/disk/by-label/nixos";
      };
    };
  };
  config = mkIf config.boot.growPartitionCustom.enable {
    # shamelessly copied from https://github.com/NixOS/nixpkgs/blob/6afb255d976f85f3359e4929abd6f5149c323a02/nixos/modules/system/boot/grow-partition.nix
    systemd.services.growpart-custom = {
      wantedBy = [ "-.mount" ];
      after = [ "-.mount" ];
      before = [ "systemd-growfs-root.service" "shutdown.target" ];
      conflicts = [ "shutdown.target" ];
      unitConfig.DefaultDependencies = false;
      serviceConfig = {
        Type = "oneshot";
        RemainAfterExit = true;
        TimeoutSec = "infinity";
        SuccessExitStatus = "0 1";
      };

      script = ''
        rootDevice="${cfg.device}"
        rootDevice="$(readlink -f "$rootDevice")"
        parentDevice="$rootDevice"
        while [ "''${parentDevice%[0-9]}" != "''${parentDevice}" ]; do
          parentDevice="''${parentDevice%[0-9]}";
        done
        partNum="''${rootDevice#''${parentDevice}}"
        if [ "''${parentDevice%[0-9]p}" != "''${parentDevice}" ] && [ -b "''${parentDevice%p}" ]; then
          parentDevice="''${parentDevice%p}"
        fi
        "${pkgs.cloud-utils.guest}/bin/growpart" "$parentDevice" "$partNum"
        "${pkgs.e2fsprogs}/bin/resize2fs" $rootDevice
      '';
    };
  };
}