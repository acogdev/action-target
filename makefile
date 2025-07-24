all: test build

build:
	go build

test:
	go test ./...

install:
	@echo "Still need to write the systemd unit file"