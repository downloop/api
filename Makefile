BINARY = api
TAG = downloop/api

all: gen api

.PHONY: api
api:
	go build -o ${BINARY} ./cmd/api

.PHONY: gen
gen: deps
	oapi-codegen -package v1 spec/api.yaml > pkg/api/v1/api.gen.go

deps:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.11.0

.PHONY: run
run: api
	./api

.PHONY: image
image:
	docker build . -t ${TAG} 
