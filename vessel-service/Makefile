link:
	rm -rf proto && ln -s ../proto proto
prot:
	protoc -I ../proto --go_out=plugins=micro:../proto ../proto/vessel/vessel.proto
build:
	protoc -I ../proto --go_out=plugins=micro:../proto ../proto/vessel/vessel.proto
	docker build -t vessel-service .
run:
	docker run -p 50052:50051 -e MICRO_SERVER_ADDRESS=:50051 -e MICRO_REGISTRY=mdns vessel-service