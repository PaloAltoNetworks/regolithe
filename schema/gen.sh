#!/bin/bash

perl -pe 's/__ATTRIBUTE__/'"$(cat rego-attribute.in)"'/g;' rego-abstract.in > rego-abstract.json
perl -pe 's/__ATTRIBUTE__/'"$(cat rego-attribute.in)"'/g;' -pe 's/__PARAMETER__/'"$(cat rego-param.in)"'/g;' rego-spec.in > rego-spec.json
perl -pe 's/__PARAMETER__/'"$(cat rego-param.in)"'/g;' rego-shared-params.in > rego-shared-params.json

rm -f ./bindata.go;
go-bindata -pkg schema .
