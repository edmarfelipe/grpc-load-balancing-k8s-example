build-proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative shared/hello.proto

build-server:
	docker build . --build-arg dir=server -t edmarfelipe/grpc-server:0.0.5

build-client:
	docker build . --build-arg dir=client -t edmarfelipe/grpc-client:0.0.5

k8s-apply:
	kubectl apply -f .k8s/