test:
	@go test -v ./...

tidy:  ## Get the dependencies
	@go mod tidy

compile: tidy  ## compiles  service code
	@go build ./...

test-race-cond:
	@go test -v -race ./...

run: tidy ## Run the scrapping server
	@go run cmd/*.go

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: compile test test-race-cond tidy run help