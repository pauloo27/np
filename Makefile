BINARY_NAME = np

.PHONY: build
build:
	go build -v -o $(BINARY_NAME) ./cmd/$(BINARY_NAME)

.PHONY: run
run: build
	./$(BINARY_NAME) 

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: lint
lint:
	go tool revive -formatter friendly -config revive.toml ./...

.PHONY: spell
spell:
	find . -name '*.go' -exec go tool misspell -error {} +

.PHONY: staticcheck
staticcheck:
	go tool staticcheck ./...

.PHONY: gosec
gosec:
	go tool gosec -tests ./...

.PHONY: inspect
inspect: spell lint staticcheck gosec

# (build but with a smaller binary)
.PHONY: dist
dist:
	go build -o $(BINARY_NAME) -ldflags="-w -s" -gcflags=all=-l -v
