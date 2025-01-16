{
  description = "bicipi";

  nixConfig = {
    extra-substituters = [
      "https://cache.garnix.io"
    ];
    extra-trusted-public-keys = [
      "cache.garnix.io:CTFPyKSLcx5RMJKfLo5EEPUObbA78b0YQ2DTCJXqr9g="
    ];
  };

  inputs = {
    # pin nixpkgs to latest working https://github.com/NixOS/nixpkgs/pull/362081#issuecomment-2596822330
    nixpkgs.url = "github:NixOS/nixpkgs?ref=7cf092925906d588daabc696d663c100f2bbacc6";
    blueprint.url = "github:numtide/blueprint";
    blueprint.inputs.nixpkgs.follows = "nixpkgs";
    nixos-hardware.url = "github:NixOS/nixos-hardware/master";
    nix-pi-loader.url  = "github:rcambrj/nix-pi-loader";
    nix-pi-loader.inputs.nixpkgs.follows = "nixpkgs";
    gomod2nix.url = "github:nix-community/gomod2nix";
    gomod2nix.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = inputs: inputs.blueprint { inherit inputs; };
}
