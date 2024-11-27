{ pkgs }:
pkgs.mkShell {
  # Add build dependencies
  packages = with pkgs; [ go libusb1 ];

  # Add environment variables
  env = { };

  # Load custom bash code
  shellHook = ''

  '';
}
