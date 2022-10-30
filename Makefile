BINARY = api
TAG = ghcr.io/downloop/api
RUNFLAGS ?= ""

all: api cli

.PHONY: api
api: gen
	go build -o ${BINARY} ./cmd/api

cli: gen
	go build -o cli ./cmd/cli

.PHONY: gen
gen: deps
	oapi-codegen -old-config-style -package v1 -generate server,spec spec/api.yaml > pkg/api/v1/api.gen.go
	oapi-codegen -old-config-style -package v1 -generate client,types spec/api.yaml > pkg/api/v1/cli.gen.go

.PHONY: deps
deps:
	go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.11.0

.PHONY: run
run: api
	./api $(RUNFLAGS) 

.PHONY: image
image:
	docker build . -t ${TAG}

.PHONY: push
push: image
	docker push ${TAG}:latest

.PHONY: test
test:
	go test -v ./pkg/api/v1/
