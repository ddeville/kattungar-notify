.PHONY: build-cli run-cli install-cli build-server run-server

build-cli:
	go build -o build/kattungar-notify-admin ./cmd/cli

run-cli:
	go run ./cmd/cli

install-cli: build-cli
	sudo mv build/kattungar-notify-admin /usr/local/bin/kattungar-notify-admin

build-server:
	docker-compose -f docker-compose.yaml build

run-server:
	docker-compose -f docker-compose.yaml up
