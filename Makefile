.PHONY: run vet fmt examples

run:
	@go run ./app

vet:
	@go vet ./...

fmt:
	@gofmt -l -w .

examples:
	@echo "hello dashboard form browser player editor chat"
