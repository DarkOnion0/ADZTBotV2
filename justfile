#!/usr/bin/env just --justfile

# Just SETTINGS (vars...)
set dotenv-load

VERSION := "latest"
export GH_TOKEN := ""
export GH_REPO := env_var_or_default("GH_REPO", "DarkOnion0/ADZTBotV2")

#Change the default just behaviour
default:
  @just --list

# Build ADZTBotV2 for all plateform
build:
    ./build.sh {{VERSION}}

# Clean the remote GHCR container registry
cleanc:
    ./delete_remote_images.sh

# Clean the binary folder
cleanb:
    rm -rf ./bin

# Lint the project files
lint:
    echo "Lint all go files"
    golangci-lint run --verbose --fix --timeout 5m .
    
    echo "Check if go.mod and go.sum are up to date"
    go mod tidy

# Format all the project files
format:
    gofmt -w .

# Shortcut to format and lint recipes
check: format lint

# Build & release ADZTBotV2, it needs GH_TOKEN to be overwritten and UNSTABLE set to unstable to publish a pre-release
release_full $UNSTABLE="stable": build
    #!/usr/bin/env bash
    if [ "${UNSTABLE}" = "unstable" ]; then
        gh release create --generate-notes --prerelease {{VERSION}} ./bin/*.zip
    else; then
        gh release create --generate-notes {{VERSION}} ./bin/*.zip
    fi

# Upload the release binary to an existing release, it needs GH_TOKEN to be overwritten
release_ci: build
    gh release upload {{VERSION}} ./bin/*.zip

# Aliases
#alias b := build
#alias cc := cleanc
#alias cb := cleanb
#alias l := lint
#alias f := format
#alias c := check

# Local Variables:
# mode: makefile
# End:
