codegen:
	@mkdir -p static
	@rm -f ./static/bindata.go;
	@go-bindata -pkg static ./templates
	@mv ./bindata.go ./static
