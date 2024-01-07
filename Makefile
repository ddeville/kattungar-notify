.PHONY: build-cli install-cli

build-cli:
	go build -o build/kattungar-notify-admin ./cmd/cli

install-cli: build-cli
	sudo mv build/kattungar-notify-admin /usr/local/bin/kattungar-notify-admin
