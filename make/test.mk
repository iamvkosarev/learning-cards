unit-test:
	@go test -v -race -cover ./internal/... -coverprofile=cover.out
	@go tool cover -html=cover.out