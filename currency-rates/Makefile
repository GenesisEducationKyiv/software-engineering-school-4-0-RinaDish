.PHONY: build
build:
	CGO_ENABLED=0 go build -a -installsuffix cgo -o $(CURDIR)/bin $(CURDIR)/cmd/main.go

.PHONY: lint 
lint:
	@echo "Running linter..."
	@golangci-lint run ./...

.PHONY: test 
test:
	@echo "Running tests..."
	go test -v ./...