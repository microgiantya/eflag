.PHONY: test
test:
	@go clean -testcache
	@go test ./... -v

.PHONY: lint
lint:
	@docker run --rm -v $(shell pwd):/repo -w /repo golangci/golangci-lint:v1.62.2-alpine golangci-lint run

.PHONY: doc
doc:
	@go install golang.org/x/tools/cmd/godoc; \
	godoc

.PHONY: fieldaligment
fieldaligment:
	@go install golang.org/x/tools/go/analysis/passes/fieldalignment/cmd/fieldalignment@latest; \
	fieldalignment -fix ./...