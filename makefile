.PHONY: all update run build install docker proto

all: proto update run

update: 
	go mod tidy \
	&& go mod vendor

run:
	go run ./pkg

build:
	cd pkg \
	&& go build

install:
	cd pkg \
	&& go install

docker:
	docker build -t lwj5/jobmaker .

proto:
	mkdir -p pkg/jobmaker \
	&& cd pkg/jobmaker \
	&& protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --proto_path ../../proto  ../../proto/jobmaker.proto
