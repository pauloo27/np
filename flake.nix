{
  description = "np - Nix project development environment manager";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-25.11";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    {
      homeManagerModules.default =
        {
          config,
          lib,
          pkgs,
          ...
        }:
        with lib;
        let
          cfg = config.programs.np;
          yamlFormat = pkgs.formats.yaml { };
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

            workspacePath = mkOption {
              type = types.nullOr types.str;
              default = null;
              description = "Path to the workspace file. Defaults to XDG_STATE_HOME/np/workspace.yaml if not set.";
            };

            tmuxBaseWindowIndex = mkOption {
              type = types.int;
              default = 1;
              description = "Base window index for tmux (usually 0 or 1).";
            };
          };

          config = mkIf cfg.enable {
            home.packages = [ cfg.package ];

            xdg.configFile."np/config.yaml".source = yamlFormat.generate "config.yaml" (
              {
                profiles_path = cfg.profilesPath;
                tmux_base_window_index = cfg.tmuxBaseWindowIndex;
              }
              // optionalAttrs (cfg.workspacePath != null) { workspace_path = cfg.workspacePath; }
            );
          };
        };
    }
    // flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = nixpkgs.legacyPackages.${system};
      in
      {
        packages.default = pkgs.buildGoModule {
          pname = "np";
          version = "0.1.0";
          src = ./.;

          vendorHash = "sha256-U7XW0TjcW/KMcnK8TuO6bMFpuG3n6n/jcTmnq+A1Vr0=";

          subPackages = [ "cmd/np" ];

          ldflags = [
            "-s"
            "-w"
          ];

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
            nil
          ];

          shellHook = ''
            export GOPATH=$HOME/go
            export PATH=$GOPATH/bin:$PATH
          '';
        };
      }
    );
}
