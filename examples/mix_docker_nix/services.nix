{ pkgs }:
let
  nativePkgs = import <nixpkgs> { };

  postgresqlPath = pkgs.lib.makeBinPath [
    pkgs.coreutils
    pkgs.gnugrep
    pkgs.postgresql_17
    pkgs.util-linux
  ];

  redisPath = pkgs.lib.makeBinPath [
    pkgs.coreutils
    pkgs.redis
  ];

  postgresqlScript = nativePkgs.writeTextFile {
    name = "start-postgres";
    destination = "/bin/start-postgres";
    executable = true;

    text = ''
      #!${pkgs.runtimeShell}
      set -euo pipefail

      export PATH=${postgresqlPath}

      : "''${PGDATA:=/var/lib/postgresql/data/pgdata}"
      : "''${POSTGRES_DB:=nozarashi}"
      : "''${POSTGRES_USER:=nozarashi}"
      : "''${POSTGRES_PASSWORD:=nozarashi}"

      export HOME=/var/lib/postgresql

      mkdir -p /etc /tmp "$HOME"
      chmod 1777 /tmp

      if ! grep -q "^postgres:" /etc/passwd 2>/dev/null; then
        printf "postgres:x:999:999:PostgreSQL:%s:/sbin/nologin\n" "$HOME" >> /etc/passwd
      fi

      if ! grep -q "^postgres:" /etc/group 2>/dev/null; then
        printf "postgres:x:999:\n" >> /etc/group
      fi

      install -d -m 0700 -o 999 -g 999 "$PGDATA"
      chown -R 999:999 "$HOME"

      if [[ ! -s "$PGDATA/PG_VERSION" ]]; then
        passwordFile="$(mktemp)"
        trap 'rm -f "$passwordFile"' EXIT

        printf "%s\n" "$POSTGRES_PASSWORD" > "$passwordFile"
        chown 999:999 "$passwordFile"
        chmod 0600 "$passwordFile"

        setpriv --reuid=999 --regid=999 --clear-groups \
          initdb \
          --pgdata="$PGDATA" \
          --username="$POSTGRES_USER" \
          --pwfile="$passwordFile" \
          --auth-host=scram-sha-256 \
          --auth-local=trust \
          --locale=C \
          --encoding=UTF8

        rm -f "$passwordFile"
        trap - EXIT

        setpriv --reuid=999 --regid=999 --clear-groups \
          pg_ctl \
          --pgdata="$PGDATA" \
          --options="-c listen_addresses= -c unix_socket_directories=/tmp" \
          --wait \
          start

        setpriv --reuid=999 --regid=999 --clear-groups \
          createdb \
          --host=/tmp \
          --username="$POSTGRES_USER" \
          "$POSTGRES_DB"

        setpriv --reuid=999 --regid=999 --clear-groups \
          pg_ctl \
          --pgdata="$PGDATA" \
          --mode=fast \
          --wait \
          stop
      fi

      # Allow authenticated connections from the project network.
      if ! grep -q "^host all all all scram-sha-256$" "$PGDATA/pg_hba.conf"; then
        printf "host all all all scram-sha-256\n" >> "$PGDATA/pg_hba.conf"
      fi

      exec setpriv --reuid=999 --regid=999 --clear-groups \
        postgres \
        -D "$PGDATA" \
        -c "listen_addresses=*" \
        -c "unix_socket_directories=/tmp"
    '';
  };

  redisScript = nativePkgs.writeTextFile {
    name = "start-redis";
    destination = "/bin/start-redis";
    executable = true;

    text = ''
      #!${pkgs.runtimeShell}
      set -euo pipefail

      export PATH=${redisPath}

      mkdir -p /data
      exec redis-server --protected-mode no --appendonly yes --dir /data
    '';
  };
in
{
  postgresql = nativePkgs.buildEnv {
    name = "postgresql-service";
    paths = [
      postgresqlScript
      pkgs.bash
      pkgs.postgresql_17
    ];
    pathsToLink = [ "/bin" ];
  };

  redis = nativePkgs.buildEnv {
    name = "redis-service";
    paths = [
      redisScript
      pkgs.bash
      pkgs.redis
    ];
    pathsToLink = [ "/bin" ];
  };
}
