#!/bin/bash
set -e

rm -rf ./schemas
mkdir ./schemas
cp ../../schema/*.json ./schemas

cd ../../cmd/rego/ || exit 1
go build
cd - || exit 1

rm -rf ./bin
mkdir ./bin
cp ../../cmd/rego/rego ./bin/rego

npm run compile
