generate:
	go generate ./...

proto:
	protoc \
		--go_out=stan --go_opt=paths=source_relative \
		--go-grpc_out=stan --go-grpc_opt=paths=source_relative \
		--proto_path=stan stan/*.proto
	
	
