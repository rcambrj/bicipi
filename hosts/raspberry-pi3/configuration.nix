{ config, flake, inputs, lib, modulesPath, pkgs, ... }: {
  imports = [
    inputs.nixos-hardware.nixosModules.raspberry-pi-3
    inputs.nix-pi-loader.nixosModules.default
    ../../modules/nixos/grow-partition.nix
    ../../modules/nixos/bicipi.nix
  ];

  system.stateVersion = "25.05";
  nixpkgs.hostPlatform = "aarch64-linux";
  boot.pi-loader.enable = true;
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
      device = "/dev/disk/by-label/ESP";
      fsType = "vfat";
    };
    "/" = {
      fsType = "tmpfs";
      options = [ "mode=0755" ];
    };
    "/mnt/root" = {
      device = "/dev/disk/by-label/nixos";
      neededForBoot = true;
      autoResize = true; # resizes filesystem to occupy whole partition
      fsType = "ext4";
    };
    "/nix" = {
      device = "/mnt/root/nix";
      neededForBoot = true;
      options = [ "defaults" "bind" ];
      depends = [ "/mnt/root" ];
    };
  };
  boot.growPartitionCustom = {
    enable = true;
    device = "/dev/disk/by-label/nixos";
  };
  system.build.image = (import "${toString modulesPath}/../lib/make-disk-image.nix" {
    inherit lib config pkgs;
    format = "raw";
    partitionTableType = "efi";
    copyChannel = false;
    diskSize = "auto";
    additionalSpace = "64M";
    bootSize = "128M";
    touchEFIVars = false;
    installBootLoader = true;
    label = "nixos";
  });
  nix.settings.trusted-users = [ "root" "bicipi" ];
  nix.settings.substituters = lib.mkForce config.nix.settings.trusted-substituters;
  nix.settings.trusted-substituters = [
    "https://cache.nixos.org/"
    "https://nix-community.cachix.org"
  ];
  nix.settings.trusted-public-keys = [
    "cache.nixos.org-1:6NCHdD59X431o0gWypbMrAURkbJ16ZPMQFGspcDShjY="
    "nix-community.cachix.org-1:mB9FSh9qf2dCimDSUo8Zy7bkq5CX+/rkCWyvRCYg3Fs="
  ];
  systemd.network.enable = true;
  networking.useDHCP = false;
  networking.useNetworkd = true;
  systemd.network = {
    networks."10-wired" = {
      matchConfig.Name = "e*";
      networkConfig = {
        # DHCP = "yes";
        DHCP = "no";
        Address = "192.168.24.24/24";
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

  services.udev.extraRules = ''
    SUBSYSTEM=="usb", ATTR{idVendor}=="3561", MODE:="0666"
  '';
  hardware.bluetooth.enable = true;
  services.bicipi = {
    enable = true;
    # debug faster:
    # extraArgs = "--calibrate=false";
  };
}