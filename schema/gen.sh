#!/bin/bash

perl -pe 's/__ATTRIBUTE__/'"$(cat rego-attribute.in)"'/g' rego-abstract.in > rego-abstract.json
perl -pe 's/__ATTRIBUTE__/'"$(cat rego-attribute.in)"'/g' rego-spec.in > rego-spec.json

rm -f ./bindata.go;
go-bindata -pkg schema .
