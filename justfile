set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

install:
    #!/usr/bin/env bash
    set -euo pipefail
    go install ./cmd/cm
    mkdir -p "${HOME}/.zfunc"
    go run ./cmd/cm completion zsh > "${HOME}/.zfunc/_cm"

    bin_dir="$(go env GOBIN)"
    if [[ -z "${bin_dir}" ]]; then
      bin_dir="$(go env GOPATH)/bin"
    fi

    printf 'Installed cm to %s and zsh completion to %s\n' "${bin_dir}/cm" "${HOME}/.zfunc/_cm"
