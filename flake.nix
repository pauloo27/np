{
  description = "np - Nix project development environment manager";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils }:
    {
      homeManagerModules.default = { config, lib, pkgs, ... }:
        with lib;
        let
          cfg = config.programs.np;
        in
        {
          options.programs.np = {
            enable = mkEnableOption "np - Nix project development environment manager";

            package = mkOption {
              type = types.package;
              default = self.packages.${pkgs.system}.default;
              defaultText = literalExpression "pkgs.np";
              description = "The np package to use.";
            };

            profilesPath = mkOption {
              type = types.str;
              default = "${config.xdg.configHome}/nix-conf/dev";
              defaultText = literalExpression ''"''${config.xdg.configHome}/nix-conf/dev"'';
              description = "Path to the directory containing nix development profiles.";
            };

            tmux = {
              windowCount = mkOption {
                type = types.int;
                default = 1;
                description = "Default number of tmux windows to create.";
              };
            };
          };

          config = mkIf cfg.enable {
            home.packages = [ cfg.package ];

            xdg.configFile."np/config.toml".text = generators.toTOML {} {
              profiles_path = cfg.profilesPath;
              tmux = {
                window_count = cfg.tmux.windowCount;
              };
            };
          };
        };
    } // flake-utils.lib.eachDefaultSystem (system:
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

          nativeBuildInputs = [ pkgs.installShellFiles ];

          postInstall = ''
            installShellCompletion --cmd np \
              --bash <($out/bin/np completion bash) \
              --fish <($out/bin/np completion fish) \
              --zsh <($out/bin/np completion zsh)
          '';

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
