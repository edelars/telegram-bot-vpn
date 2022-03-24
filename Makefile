swagger:
	@swagger generate server -f=swagger/swagger.yaml -t internal/api --exclude-main

proto:
	@protoc --go_out=:. --proto_path=${GOPATH}/src --proto_path=internal/proto internal/proto/*.proto