run:
	@go run ./main.go

gen:
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:proto

rm:
	rm proto/*.go