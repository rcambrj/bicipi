{ flake, perSystem, pkgs }:
pkgs.mkShell {
  # Add build dependencies
  packages = with pkgs; [
    go
    gopls
    perSystem.gomod2nix.gomod2nix
    (perSystem.gomod2nix.mkGoEnv { pwd = ./.; })
    usbutils
    pkg-config # required for libusb1
    libusb1
    # flake.packages.${pkgs.stdenv.hostPlatform.system}.bicipi
  ];

  # Add environment variables
  env = { };

  # Load custom bash code
  shellHook = ''

  '';
}
