all: test build

build:
	go build

test:
	go test ./...

install:
	mkdir -p /usr/bin/action-target
	mkdir -p /etc/action-target
	cp action-target /usr/bin/action-target/
	cp config.toml /etc/action-target/
	cp action-target.service /usr/lib/systemd/system/
	systemctl daemon-reload
	systemctl enable action-target.service
	systemctl start action-target.service