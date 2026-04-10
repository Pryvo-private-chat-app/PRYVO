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
	gsettings-desktop-schemas
      ];

      shellHook = ''
	echo "Hello There!"
        echo "Go version: $(go version)"
        echo "Node version $(node -v)"
        
        export XDG_DATA_DIRS=${pkgs.gsettings-desktop-schemas}/share/gsettings-schemas/${pkgs.gsettings-desktop-schemas.name}:${pkgs.gtk3}/share/gsettings-schemas/${pkgs.gtk3.name}:$XDG_DATA_DIRS
      '';
    };
  };
} 
