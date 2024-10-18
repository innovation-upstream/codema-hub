{ pkgs ? import <nixpkgs> {}, ... }:

pkgs.mkShell {
  buildInputs = with pkgs;
  let
  in
  [
    tilt
    go_1_23
    kube3d
    kubectl
  ];
}

