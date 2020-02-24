SHELL:=/bin/bash -o pipefail

CGO_ENABLED=0

#########################################
#	Testing and building the binary
#########################################

.PHONY: deps
deps:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.16.0
	go get -u github.com/jstemmer/go-junit-report
	go get github.com/deepmap/oapi-codegen/cmd/oapi-codegen@master

.PHONY: lint
lint:
	golangci-lint run ./...

.PNONY: test
test:
	go test -covermode=atomic -coverprofile=./coverage.txt -v ./... 2>&1 | tee ./test.txt
	cat test.txt | go-junit-report > ./report.xml
	go tool cover -func=./coverage.txt

.PHONY: openapi
openapi:
	mkdir -p api
	@echo Generating API clients
	oapi-codegen -generate "types,client,spec" --package api ../notifier/api/spec.yaml > api/api.gen.go
