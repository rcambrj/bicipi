{ flake, perSystem, pkgs }:
pkgs.mkShell {
  # Add build dependencies
  packages = with pkgs; [
    go
    gopls
    perSystem.gomod2nix.gomod2nix
    usbutils
    libusb1
    pkg-config # required for gousb
    flake.packages.${pkgs.stdenv.hostPlatform.system}.bicipi
  ];

  # Add environment variables
  env = { };

  # Load custom bash code
  shellHook = ''

  '';
}
