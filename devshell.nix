{ pkgs }:
pkgs.mkShell {
  # Add build dependencies
  packages = [ "go" "libusb1" ];

  # Add environment variables
  env = { };

  # Load custom bash code
  shellHook = ''

  '';
}
