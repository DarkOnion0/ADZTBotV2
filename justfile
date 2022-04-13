# Change the default just behaviour
default:
  @just --list

# Build ADZTBotV2 for all plateform
build VERSION="":
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

# Aliases
#alias b := build
#alias cc := cleanc
#alias cb := cleanb
#alias l := lint
#alias f := format
#alias c := check

# TODO: make a release command and update the CI
