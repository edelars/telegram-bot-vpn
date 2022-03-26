APP = vpn-backend
NAMESPACE = backend
BUILD_DIR = build
CGO_ENABLED = 0
GOTAGS ?= musl
GOOS ?= linux



swagger:
	@swagger generate server -f=swagger/swagger.yaml -t internal/api --exclude-main

proto:
	@protoc --go_out=:. --proto_path=${GOPATH}/src --proto_path=internal/proto internal/proto/*.proto

build:
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -mod=vendor -tags $(GOTAGS) -ldflags '-s -w -extldflags "-static"' -o ./${BUILD_DIR}/${APP} ./cmd/main.go

install:
	sudo cp ./${BUILD_DIR}/${APP} /usr/local/sbin/${APP}