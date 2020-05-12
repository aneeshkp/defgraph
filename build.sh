#!/usr/bin/env bash
set -e

DEST_DIR="bin"

if [ ! -d ${DEST_DIR} ]; then
        mkdir ${DEST_DIR}
fi


        # build with go modules
        export GO111MODULE=on

        echo "Building dep graph"
        go build -o ${DEST_DIR}/depgraph  -ldflags "${LDFLAGS}"
f
# go install ./...