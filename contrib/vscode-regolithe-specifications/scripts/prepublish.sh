#!/bin/bash
set -e

cd ../../schema || exit 1
./gen.sh || exit 1
cd - || exit 1

rm -rf ./schemas
mkdir ./schemas
cp ../../schema/*.json ./schemas || exit 1

cd ../../cmd/rego/ || exit 1
go build || exit 1
cd - || exit 1

rm -rf ./bin
mkdir ./bin
cp ../../cmd/rego/rego ./bin/rego || exit 1

npm run compile || exit 1
