#!/bin/bash
set -e

(
    cd ../../schema
    ./gen.sh
)

rm -rf ./schemas
mkdir ./schemas
cp ../../schema/*.json ./schemas

(
    cd ../../cmd/rego/
    go generate
    go install
    go build
)

rm -rf ./bin
mkdir ./bin
cp ../../cmd/rego/rego ./bin/rego

npm run compile
