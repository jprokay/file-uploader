{ pkgs, lib, config, inputs, ... }:

let
  pkgs-unstable = import inputs.nixpkgs-unstable { system = pkgs.stdenv.system; };
in {
  # https://devenv.sh/basics/
  env.GREET = "devenv";

  # https://devenv.sh/packages/
  packages = [ 
    pkgs-unstable.go
    pkgs-unstable.gopls
    pkgs.git
    pkgs.dbmate
    pkgs.corepack
    pkgs.prettierd
    pkgs.eslint_d
    pkgs.postgresql_16
  ];

  # https://devenv.sh/services/
  services.postgres = {
    enable = true;
    package = pkgs.postgresql_16;
    listen_addresses = "127.0.0.1";
    port = 5432;
    initialDatabases = [
      {
        name = "db";
      }
      {
        name = "testdb";
      }
    ];
    initialScript = "
      CREATE USER dbmate WITH SUPERUSER PASSWORD 'dbmate';
    ";
  };

  env.DATABASE_URL = "postgresql://dbmate:dbmate@127.0.0.1:5432/db?sslmode=disable";
  env.DBMATE_MIGRATIONS_DIR = "./db/migrations";
  env.DBMATE_SCHEMA_FILE = "./db/schema.sql";

  env.BACKEND_ALLOW_ORIGIN = "http://localhost:3000";
  env.PORT = "8080"; 
  env.BACKEND_URL = "http://localhost:8080";
  env.NEXT_PUBLIC_BACKEND_URL = "http://localhost:8080";

  env.GOPLS = "${pkgs-unstable.gopls}/bin/gopls";

  # https://devenv.sh/languages/
  languages.typescript = {
    enable = true;
  };
  languages.go = {
    enable = true;
    package = pkgs-unstable.go;
  };

  languages.python = {
    enable = true;
    uv.enable = true;
  };

  # https://devenv.sh/pre-commit-hooks/
  # pre-commit.hooks.shellcheck.enable = true;

  # https://devenv.sh/processes/
  processes = {
    frontend.exec = "cd frontend && npm run dev";
    backend.exec = "cd backend && go run server.go"; 
  };

  # See full reference at https://devenv.sh/reference/options/
}
