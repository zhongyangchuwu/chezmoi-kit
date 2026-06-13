set shell := ["bash", "-eu", "-o", "pipefail", "-c"]

install:
    go install ./cmd/cm
    mkdir -p "$HOME/.zfunc"
    go run ./cmd/cm completion zsh > "$HOME/.zfunc/_cm"
    echo "Installed cm to $(go env GOBIN || go env GOPATH)/bin and zsh completion to ${HOME}/.zfunc/_cm"
