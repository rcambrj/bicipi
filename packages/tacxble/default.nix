{ perSystem, ... }:
perSystem.gomod2nix.buildGoApplication {
  pname = "tacxble";
  version = "0.1";
  src = ./.;
  modules = ../../gomod2nix.toml;
}