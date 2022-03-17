.PHONY: update master release setup update_master update_release build clean

setup:
	git config --global --add url."git@gitlab.com:".insteadOf "https://gitlab.com/"

clean:
	rm -rf vendor/
	go mod vendor

update:
	-GOFLAGS="" go get -u all

build:
	go build ./...
	go mod tidy

update_release:
	GOFLAGS="" go get -u gitlab.com/xx_network/primitives@release

update_master:
	GOFLAGS="" go get -u gitlab.com/xx_network/primitives@master

master: update_master clean build

release: update_release clean build
