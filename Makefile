.PHONY: test
test:
	@go clean -testcache
	@go test ./... -v

.PHONY: lint
lint:
	@docker run --rm -v $(shell pwd):/repo -w /repo golangci/golangci-lint:v1.62.2-alpine golangci-lint run
