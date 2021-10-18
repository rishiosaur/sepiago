{
  description = "A minimal functional programming language";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }: flake-utils.lib.eachDefaultSystem (system: let
    pkgs = nixpkgs.legacyPackages.${system};
  in rec {
    packages.sepia = pkgs.buildGoModule {
      name = "sepia";
      src = ./.;
      vendorSha256 = null;
      meta = with pkgs.lib; {
        description = "A minimal functional programming language";
        homepage = "https://github.com/rishiosaur/sepiago";
        license = licenses.asl20;
        platforms = platforms.linux ++ platforms.darwin;
      };
    };
    defaultPackage = packages.sepia;
  });
}
