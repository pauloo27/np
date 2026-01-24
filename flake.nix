{
  description = "np - Nix project development environment manager";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "np";
          version = "0.1.0";
          src = ./.;

          vendorHash = "sha256-n/12IYu5NOa7cY6YvsDsBVQwBjzrerHe0tNJlO4Kjlo=";

          subPackages = [ "cmd/np" ];

          ldflags = [ "-s" "-w" ];

          meta = with pkgs.lib; {
            description = "Nix project development environment manager";
            homepage = "https://code.db.cafe/pauloo27/np";
            license = licenses.mit;
            maintainers = [ ];
          };
        };

        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            go
            gopls
          ];

          shellHook = ''
            export GOPATH=$HOME/go
            export PATH=$GOPATH/bin:$PATH
          '';
        };
      }
    );
}
