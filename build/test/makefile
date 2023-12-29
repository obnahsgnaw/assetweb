
.PHONY: test_help
test_help:
	@echo "test options: <command>"
	@echo "    clean  : Clear all build binary of this project"
	@echo "    test   : Test this project"
	@echo "    check  : Check this project"
	@echo "    cover  : Cover this project"
	@echo "    lint   : Lint this project"

.PHONY: clean
clean:
	@echo "Cleaning..."
	@go clean
	@rm -rf ${OUT_NAME}/*
	@echo "Done"

.PHONY: test
test:
	@echo "Testing..."
	@go test
	@echo "Done"

.PHONY: check
check:
	@echo "Checking..."
	@go fmt ./
	@go vet ./
	@echo "Done"

.PHONY: cover
cover:
	@echo "Covering..."
	@go test -coverprofile cover.out
	@go tool cover -html=cover.out
	@echo "Done"

.PHONY: lint
lint:
	@echo "Linting..."
	@golangci-lint run
	@echo "Done"