{ perSystem, pkgs, ... }:
perSystem.gomod2nix.buildGoApplication {
  pname = "bicipi";
  version = "0.1";
  pwd = ./..;
  src = ./..;
  modules = ./../gomod2nix.toml;
  go = pkgs.go;
  nativeBuildInputs = with pkgs; [
    pkg-config
  ];
  buildInputs = with pkgs; [
    libusb1
  ];
}