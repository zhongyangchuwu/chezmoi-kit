set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

version := env_var_or_default("VERSION", "dev")
ldflags := "-X github.com/zhongyangchuwu/cm/internal/build.Version=" + version

install:
    #!/usr/bin/env bash
    set -euo pipefail
    go install -ldflags "{{ldflags}}" ./cmd/cm
    mkdir -p "${HOME}/.zfunc"
    go run ./cmd/cm completion zsh > "${HOME}/.zfunc/_cm"

    bin_dir="$(go env GOBIN)"
    if [[ -z "${bin_dir}" ]]; then
      bin_dir="$(go env GOPATH)/bin"
    fi

    printf 'Installed cm %s to %s and zsh completion to %s\n' "{{version}}" "${bin_dir}/cm" "${HOME}/.zfunc/_cm"

build-release:
    mkdir -p dist
    go build -trimpath -ldflags "-s -w {{ldflags}}" -o dist/cm ./cmd/cm
