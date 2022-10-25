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
build: install
    ./build.sh {{VERSION}}

# Clean the remote GHCR container registry
cleanc:
    ./delete_remote_images.sh

# Clean the binary folder
cleanb:
    rm -rf ./bin

# Lint the project files
lint: install
    @echo -e "\nLint all go files"
    golangci-lint run --verbose --fix --timeout 5m .

# Format all the project files
format:
    @echo -e "\nFormat go code"
    gofmt -w .

    @echo -e "\nFormat other code with prettier (yaml, md...)"
    prettier -w .

# Check the go.mod and the go.sum files
check: install format lint
    @echo -e "\nVerify dependencies have expected content"
    go mod verify
    
    @echo -e "\nCheck if go.mod and go.sum are up to date"
    go mod tidy

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

# The command to run to dev ADZTBotV2
dev: format lint
    @echo -e "\nRun ADZTBotV2"
    env ADZTBOTV2_V1_BEFORE=true go run -ldflags="-X 'github.com/DarkOnion0/ADZTBotV2/config.RawVersion={{VERSION}}'" main.go -db $DB -url $URL -chanm $CHANM -chanv $CHANV -token $TOKEN -admin $ADMIN -debug $DEBUG -timer $TIMER

# Run the prerequisites to install all the missing deps that nix can't cover
install:
    go mod download

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
