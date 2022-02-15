.PHONY:
.SILENT:

build:
	go build -o ./.bin/calloriecounter ./cmd/bot/main.go

run: build
	./.bin/calloriecounter