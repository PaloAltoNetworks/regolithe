codegen:
	@rm -f ./schema/bindata.go;
	@go-bindata -pkg schema ./schema
	@mv ./bindata.go ./schema
