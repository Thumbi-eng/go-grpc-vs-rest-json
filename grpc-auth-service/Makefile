.PHONY: proto tidy run-server run-client clean

PROTO_DIR=proto
PB_DIR=proto/pb

proto:
	@mkdir -p $(PB_DIR)
	protoc --proto_path=$(PROTO_DIR) \
		--go_out=$(PB_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

tidy:
	go mod tidy

run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/client/main.go

clean:
	rm -f $(PB_DIR)/*.pb.go
