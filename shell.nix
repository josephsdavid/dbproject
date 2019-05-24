let
  pkgs = import <nixpkgs> {};
in
  pkgs.mkShell {
    buildInputs = with pkgs; [
      texlive.combined.scheme-full
      go
      pkg-config
      openssl_1_1
      cmake
      python37
      R
      neo4j
      mongodb
      mysql
    ];
  }
