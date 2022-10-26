BINARY = api
TAG = ghcr.io/downloop/api

all: api

.PHONY: api
api: gen
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

.PHONY: push
push: image
	docker push ${TAG}:latest

.PHONY: test
	act -j build
