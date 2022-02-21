.PHONY:
.SILENT:

build:
	go build -o ./.bin/calloriecounter ./cmd/bot/main.go

run: build
	./.bin/calloriecounter

build_image:
	docker build -t calorie_counter_bot:v.0.1

start_container:
	docker run --name calorie_counter_bot --env-file .env