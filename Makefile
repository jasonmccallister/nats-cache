.PHONY: proto nats-server
PUB_KEY ?= f0e061f5981a2b0bb3d6a5b7e1e7c557d4c8ec6fdc9eef98c37b6b2983a0b912

build:
	go build -o bin/nats-cache cmd/nats-cache/main.go
run: build
	./bin/nats-cache -public-key=$(PUB_KEY) -log-level=info -log-format=text -nats-port 4222 -nats-http-port 8222

generate:
	buf generate
nats-server:
	nats-server -js -m 8222 -p 4222
