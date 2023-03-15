#!/bin/sh

set -e

REF="refs/changes/97/476097/3"
SRC="https://go.googlesource.com/go"

really_build_go() {
    (
        cd go/src
        ./make.bash
    )
}

build_go() {
    if [ ! -f "./go/bin/go" ]; then
        really_build_go
        return
    fi
    goVersion="$(./go/bin/go version | awk '{print $4}' | awk -F'-' '{print $2}')"
    if [ ! "${goVersion}" = "${commit}" ]; then
        echo "Go version does not match expected, rebuilding: ${goVersion} != ${commit}" >&2
        really_build_go
    fi
}

if [ ! "${USE_SYSTEM_GO}" = "1" ]; then
    if [ ! -d "./go/.git" ]; then
        mkdir go
        git -C go init .
    fi

    git -C go fetch --depth=1 "${SRC}" "${REF}" && git -C go checkout FETCH_HEAD
    commit="$(git -C go rev-parse --short HEAD)"

    build_go

    export GOROOT="${PWD}/go"
    export PATH="${PWD}/go/bin:${PATH}"
fi

go run .
