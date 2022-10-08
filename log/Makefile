.PHONY: lint
lint:
	/home/sam/go/bin/golangci-lint run

.PHONY: test
test:
	@go test ./... -coverpkg ./... -race -covermode=atomic -coverprofile=coverage.out

