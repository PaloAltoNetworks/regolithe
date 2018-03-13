codegen:
	@mkdir -p ./static
	@rm -f ./static/bindata.go;
	@cd ../.. && go-bindata -pkg static schema
	@mv ../../bindata.go ./static
