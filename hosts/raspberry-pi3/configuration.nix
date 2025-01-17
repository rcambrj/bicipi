{ config, flake, inputs, lib, modulesPath, pkgs, ... }:
with lib;
let
  debug = false;
in {
  imports = [
    inputs.nixos-hardware.nixosModules.raspberry-pi-3
    inputs.nix-pi-loader.nixosModules.default
    ../../modules/nixos/grow-partition.nix
    ../../modules/nixos/bicipi.nix
  ];

  system.stateVersion = "25.05";
  nixpkgs.hostPlatform = "aarch64-linux";
  boot.pi-loader = {
    enable = true;
    bootMode = "direct";
  };
  boot.consoleLogLevel = lib.mkDefault 7;
  nix.extraOptions = ''
    experimental-features = nix-command flakes
  '';
  nixpkgs.config.allowUnfree = true;
  services.journald.extraConfig = ''
    Storage=volatile
  '';
  networking.firewall.enable = true;
  fileSystems = {
    "/boot" = {
      device = "/dev/disk/by-label/BOOT";
      fsType = "vfat";
      options = mkIf (!debug) [ "ro" ];
    };
    "/" = {
      fsType = "tmpfs";
      options = [ "mode=0755" ];
    };
    "/mnt/root" = {
      device = "/dev/disk/by-label/nixos";
      neededForBoot = true;
      autoResize = !debug; # resizes filesystem to occupy whole partition
      fsType = "ext4";
      options = mkIf (!debug) [ "ro" ];
    };
    "/nix" = {
      device = "/mnt/root/nix";
      neededForBoot = true;
      options = [ "defaults" "bind" ];
      depends = [ "/mnt/root" ];
    };
  };
  boot.growPartitionCustom = {
    enable = !debug;
    device = "/dev/disk/by-label/nixos";
  };
  system.build.image = (import inputs.nix-pi-loader.nixosModules.make-disk-image {
    inherit lib config pkgs;
    format = "raw";
    partitionTableType = "legacy+boot";
    copyChannel = false;
    diskSize = "auto";
    additionalSpace = "64M";
    bootSize = "128M";
    touchEFIVars = false;
    installBootLoader = true;
    label = "nixos";
  });
  nix.settings.trusted-users = [ "bicipi" ];
  systemd.network.enable = true;
  networking.useDHCP = false;
  networking.useNetworkd = true;
  systemd.network = {
    networks."10-wired" = {
      matchConfig.Name = "e*";
      networkConfig = if debug then {
        DHCP = "yes";
      } else {
        # with DHCP this machine is likely to gain a route to the Internet
        # prevent it by disabling DHCP, but preserve emergency SSH via static IP
        DHCP = "no";
        Address = "192.168.24.24/24";
      };
    };
  };
  services.auto-cpufreq = {
    enable = true;
    settings = {
      charger = {
        governor = "powersave";
        energy_performance_preference = "power";
        turbo = "never";
      };
    };
  };
  zramSwap = {
    enable = true;
    memoryPercent = 50;
  };
  services.openssh.enable = true;
  users.users.bicipi = {
    isNormalUser = true;
    uid = 1000;
    group = "bicipi";
    home = "/home/bicipi";
    extraGroups = [
      "wheel"
      "dialout" # for serial permission
    ];
    initialPassword = "password";
  };
  users.groups.bicipi = {
    gid = 1000;
  };

  environment.systemPackages = with pkgs; [ libraspberrypi ];

  services.udev.extraRules = ''
    SUBSYSTEM=="usb", ATTR{idVendor}=="3561", MODE:="0666"
  '';
  hardware.bluetooth.enable = true;
  services.bicipi = if debug then {
    enable = true;
    extraArgs = "--calibrate=false";
  } else {
    enable = true;
  };
}