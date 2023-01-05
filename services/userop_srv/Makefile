run:
	@go run ./main.go

gen:
	@if [ -f proto/*.go ] ; then rm proto/*.go ; fi
	protoc --proto_path=proto proto/*.proto --go_out=plugins=grpc:proto