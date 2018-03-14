#!/bin/bash
set -e

rm -rf ./schemas
mkdir ./schemas
cp ../../../schema/regolithe.json ./schemas

cd ../../../cmd/rego/ || exit 1
make codegen
go build
cd - || exit 1

rm -rf ./bin
mkdir ./bin
cp ../../../cmd/rego/rego ./bin/rego

npm run compile
