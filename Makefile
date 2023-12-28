.PHONY: proto nats-server

generate:
	buf generate proto
nats-server:
	nats-server -js -m 8222 -p 4222
