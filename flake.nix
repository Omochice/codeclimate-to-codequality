{
  description = "A Go command-line tool that converts CodeClimate JSON to GitLab Code Quality format";

  nixConfig = {
    extra-substituters = [
      "https://omochice.cachix.org"
    ];
    extra-trusted-public-keys = [
      "omochice.cachix.org-1:d+cdfbGVPgtxxdGSkGf3hhaCdfziMtZ6FSHUWxwUTo8="
    ];
  };

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    nur-packages = {
      url = "github:Omochice/nur-packages";
      inputs.nixpkgs.follows = "nixpkgs";
    };
    git-hooks = {
      url = "github:cachix/git-hooks.nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    {
      self,
      nixpkgs,
      treefmt-nix,
      flake-utils,
      nur-packages,
      git-hooks,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [
            nur-packages.overlays.default
          ];
        };
        treefmt = treefmt-nix.lib.evalModule pkgs (
          { ... }:
          let
            rumdlConfig = (pkgs.formats.toml { }).generate "rumdl.toml" {
              # keep-sorted start
              MD004.style = "dash";
              MD007.indent = 4;
              MD007.style = "fixed";
              MD041.enabled = false;
              MD049.style = "underscore";
              MD050.style = "asterisk";
              MD055.style = "leading-and-trailing";
              MD060.enabled = true;
              MD060.style = "aligned";
              MD077.enabled = false;
              global.line_length = 0;
              # keep-sorted end
            };
          in
          {
            settings.global.excludes = [
              "CHANGELOG.md"
            ];
            settings.formatter.rumdl-format.options = [
              "--config"
              (toString rumdlConfig)
            ];
            programs = {
              # keep-sorted start block=yes
              formatjson5 = {
                enable = true;
                indent = 2;
              };
              gofmt.enable = true;
              goimports.enable = true;
              keep-sorted.enable = true;
              nixfmt.enable = true;
              rumdl-format.enable = true;
              toml-sort.enable = true;
              yamlfmt = {
                enable = true;
                settings = {
                  formatter = {
                    type = "basic";
                    retain_line_breaks_single = true;
                    scan_folded_as_literal = true;
                  };
                };
              };
              # keep-sorted end
            };
          }
        );
        pre-commit-check = git-hooks.lib.${system}.run {
          src = ./.;
          package = pkgs.prek;
          hooks = {
            # Claude Code sets CLAUDECODE; humans are covered by the pre-push hook
            gitleaks-commit = {
              enable = true;
              name = "gitleaks (claude commit)";
              entry = pkgs.lib.getExe (
                pkgs.writeShellApplication {
                  name = "gitleaks-when-claude";
                  runtimeInputs = [ pkgs.gitleaks ];
                  text = ''
                    if [ -z "''${CLAUDECODE:-}" ]; then
                      exit 0
                    fi
                    gitleaks git --pre-commit --staged --no-banner --redact
                  '';
                }
              );
              pass_filenames = false;
              stages = [ "pre-commit" ];
            };
            gitleaks-push = {
              enable = true;
              name = "gitleaks";
              entry = "${pkgs.lib.getExe pkgs.gitleaks} git --no-banner --redact";
              pass_filenames = false;
              stages = [ "pre-push" ];
            };
            treefmt = {
              enable = true;
              packageOverrides.treefmt = treefmt.config.build.wrapper;
              stages = [ "pre-push" ];
            };
          };
        };
        version = pkgs.lib.pipe ./.github/release-please-manifest.json [
          builtins.readFile
          builtins.fromJSON
          (builtins.getAttr ".")
        ];
        devPackages = rec {
          # keep-sorted start block=yes
          actions = with pkgs; [
            actionlint
            ghalint
            zizmor
          ];
          # keep-sorted end
          default = [
            pkgs.go_1_26
            pkgs.goreleaser
            treefmt.config.build.wrapper
          ]
          ++ actions;
        };
      in
      {
        # keep-sorted start block=yes
        checks = {
          # keep-sorted start
          actions =
            pkgs.runCommand "check-actions"
              {
                buildInputs = devPackages.actions;
                src = self;
              }
              ''
                cd $src
                actionlint .github/workflows/*.yaml
                ghalint run
                zizmor .github/workflows .github/actions
                touch $out
              '';
          formatting = treefmt.config.build.check self;
          pre-commit = pre-commit-check;
          renovate =
            pkgs.runCommand "validate-renovate-config"
              {
                buildInputs = with pkgs; [
                  renovate
                ];
                src = self;
              }
              ''
                cd $src
                renovate-config-validator --strict renovate.json5
                touch $out
              '';
          # keep-sorted end
        };
        devShells = pkgs.lib.pipe devPackages [
          (pkgs.lib.attrsets.mapAttrs (
            name: buildInputs:
            pkgs.mkShell {
              buildInputs = buildInputs ++ pre-commit-check.enabledPackages;
              inherit (pre-commit-check) shellHook;
            }
          ))
        ];
        formatter = treefmt.config.build.wrapper;
        packages = {
          default = pkgs.buildGo126Module {
            #keep-sorted start block=yes
            env.CGO_ENABLED = 0;
            ldflags = [
              "-s"
              "-w"
              "-X main.version=${version}"
            ];
            meta.description = "A Go command-line tool that converts CodeClimate JSON to GitLab Code Quality format";
            meta.homepage = "https://github.com/Omochice/codeclimate-to-codequality";
            meta.license = pkgs.lib.licenses.zlib;
            pname = "codeclimate-to-codequality";
            src = ./.;
            vendorHash = "sha256-W6XVd68MS0ungMgam8jefYMVhyiN6/DB+bliFzs2rdk=";
            version = version;
            #keep-sorted end
          };
        };
        # keep-sorted end
      }
    );
}
