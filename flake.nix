{
  description = "ssm-parent";

  inputs.nixpkgs.url = "nixpkgs/nixos-23.05";
  inputs.devshell.url = "github:numtide/devshell";
  inputs.flake-utils.url = "github:numtide/flake-utils";

  outputs = { self, nixpkgs, devshell, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        lastModifiedDate = self.lastModifiedDate or self.lastModified or "19700101";
        # Generate a user-friendly version number.
        version = builtins.substring 0 8 lastModifiedDate;

        vendorSha256 = "sha256-QbCCaOJVwToT1ViRo2CrNiL7yArfyK1MXvo9IcDNpJM="; # pkgs.lib.fakeSha256;

        pkgs = import nixpkgs {
          inherit system;
          overlays = [ devshell.overlays.default ];
        };
      in
      {
        packages.default = self.packages."${system}".ssm-parent;

        packages.ssm-parent-static =
          pkgs.buildGoModule {
            pname = "ssm-parent";
            inherit version;
            inherit vendorSha256;
            src = ./.;

            nativeBuildInputs = [ pkgs.musl ];
            buildInputs = [ ];
            pkgConfigModules = [ ];

            CGO_ENABLED = 0;
            ldflags = [
              "-linkmode external"
              "-extldflags '-static -L${pkgs.musl}/lib'"
            ];

            meta = with pkgs.lib; {
              license = [ licenses.asl20 ];
              platforms = platforms.linux;
              description = "ssm-parent fully static build";
            };
          };

        packages.ssm-parent =
          pkgs.buildGoModule {
            pname = "ssm-parent";
            inherit version;
            inherit vendorSha256;
            src = ./.;

            meta = with pkgs.lib; {
              license = [ licenses.asl20 ];
              platforms = platforms.all;
            };
          };

        devShell =
          pkgs.devshell.mkShell {
            commands = [
              { package = pkgs.go; }
              {
                package = pkgs.delve;
                name = "dlv";
              }
            ];

            packages = with pkgs; [
              gopls
              gotools
              go-tools
              rnix-lsp
            ];
          };
      });
}
