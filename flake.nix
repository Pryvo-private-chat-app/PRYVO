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
      nativeBuildInputs = with pkgs; [
	pkg-config
      ];

      buildInputs = with pkgs; [
	go
	gopls
	wails
	nodejs
	gtk3
	webkitgtk_4_1
      ];

      shellHook = ''
	echo "Hello There!"
        echo "Go version: $(go version)"
	echo "Node version $(node -v)"
      '';
    };
  };
} 
