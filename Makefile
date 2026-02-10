.PHONY: gen-pdf test-coverage

gen-pdf:
	npx --yes md-to-pdf $(FILE)

test-coverage:
	go test -v -coverpkg=./internal/...,./pkg/... ./tests/unit/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	rm coverage.out
	@echo "Coverage report generated at coverage.html"
