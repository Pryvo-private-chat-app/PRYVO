{
  description = "Direnv for PRYVO";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
  };
  
  outputs = { self, nixpkgs }:
  let
    system = "x86_64-linux";
    pkgs = nixpkgs.legacyPackages.${system};
  in {
    devShells.${system}.default = pkgs.mkShell {
      buildInputs = with pkgs; [
	go
	gopls
      ];

      shellHook = ''
	echo "Hello There!"
        echo "Go version: $(go version)"
      '';
    };
  };
} 
