{
  description = "bicipi";

  # Add all your dependencies here
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs?ref=nixos-unstable";
    blueprint.url = "github:numtide/blueprint";
    blueprint.inputs.nixpkgs.follows = "nixpkgs";
    nixos-hardware.url = "github:NixOS/nixos-hardware/master";
    nix-pi-loader.url  = "github:rcambrj/nix-pi-loader";
    nix-pi-loader.inputs.nixpkgs.follows = "nixpkgs";
    gomod2nix.url = "github:nix-community/gomod2nix";
  };

  # Load the blueprint
  outputs = inputs: inputs.blueprint { inherit inputs; };
}
