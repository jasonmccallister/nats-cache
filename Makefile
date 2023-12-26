proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=internal/gen/proto/cache --go-grpc_opt=paths=source_relative protob/*.proto
nats-server:
	nats-server -js -m 8222 -p 4222
