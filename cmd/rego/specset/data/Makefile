MAKEFLAGS += --warn-undefined-variables
SHELL := /bin/bash -o pipefail

default: format codegen

.PHONY:codegen
codegen:
	@ elegen folder -d specs -o codegen -g openapi3 || exit 1
	@ rm -rf openapi3
	@ elegen folder -d specs -o codegen || exit 1
	@ mv custom_validations.go custom_validations.go.keep
	@ mv custom_validations_test.go custom_validations_test.go.keep
	@ rm -rf ./*.go
	@ mv custom_validations.go.keep custom_validations.go
	@ mv custom_validations_test.go.keep custom_validations_test.go
	@ mv codegen/elemental/*.go ./
	@ mv codegen/openapi3 openapi3
	@ rm -rf codegen doc
	@ mkdir -p doc
	@ data=$$(rego doc -d specs || exit 1) && echo -e "$${data}" > doc/documentation.md

format: format-specs format-type format-validation format-parameter
format-specs:
	@ for f in specs/*.abs; do \
		rego format < $$f > $$f.formatted && \
		mv $$f.formatted $$f; \
	done
	@ for f in specs/*.spec; do \
		rego format < $$f > $$f.formatted && \
		mv $$f.formatted $$f; \
	done

format-type: target = "specs/_type.mapping"
format-type:
	@ rego format -m typemapping < $(target) > $(target).formatted
	@ mv $(target).formatted $(target)

format-validation: target = "specs/_validation.mapping"
format-validation:
	@ rego format -m validationmapping < $(target) > $(target).formatted
	@ mv $(target).formatted $(target)

format-parameter: target = "specs/_parameter.mapping"
format-parameter:
	@ rego format -m parametermapping < $(target) > $(target).formatted
	@ mv $(target).formatted $(target)
