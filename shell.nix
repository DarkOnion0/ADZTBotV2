{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  shellHook = ''
      # Go command
      go mod verify
      go mod tidy
      go mod download

      # Alias
      alias build="bash ./build.sh"
      alias cleanc="bash ./delete_remote_images.sh"
      alias cleanb="rm -rf ./bin"

      # Welcome script
      echo -e "\n$(tput bold)Welcome in the nix-shell for ADZTBotV2$(tput sgr0)"
      
      echo -e "\nList of custom command (only supputed in BASH) :"
      echo -e "=================================================="
      echo -e "- build    | run the compile script, accept 1 argument that will be used as the version tag"
      echo -e "- cleanc   | clean the remote untagged ghcr images"
      echo -e "- cleanb   | clean the build folder"
      echo -e "=================================================="
    '';

    # nativeBuildInputs is usually what you want -- tools you need to run
    nativeBuildInputs = [
      pkgs.go
      pkgs.zip
      pkgs.unzip
      pkgs.curl
      pkgs.jq
     ];
}
