coverage:
	go tool cover -func=coverage.out
generate:
	$(ENV_VARS) go generate ./...
test:
	go test -v -count=1 -coverprofile coverage.out -race ./...
tidy:
	go mod tidy