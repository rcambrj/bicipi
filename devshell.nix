{ pkgs }:
pkgs.mkShell {
  # Add build dependencies
  packages = with pkgs; [
    go
    gopls
    usbutils
    libusb1
    pkg-config # required for gousb
  ];

  # Add environment variables
  env = { };

  # Load custom bash code
  shellHook = ''

  '';
}
