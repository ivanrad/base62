GOCMD := go
GOTEST := $(GOCMD) test
GOBUILD := $(GOCMD) build

.PHONY: build test fuzz cover vet staticcheck gofmt clean check lint ci bench gofumpt

build:
	$(GOBUILD)

test:
	$(GOTEST) ./...

fuzz:
	$(GOTEST) -fuzz ./...

cover:
	$(GOTEST) -coverprofile base62.out ./...

vet:
	$(GOCMD) vet ./...

staticcheck:
#	$(GOCMD) run honnef.co/go/tools/cmd/staticcheck@latest -- $$(go list ./...)
	staticcheck ./...

gofmt:
	@echo "gofmt -l ./"
	@test -z "$$(gofmt -l ./ | tee /dev/stderr)"

gofumpt:
	@echo "gofumpt -l ./"
	@test -z "$$(gofumpt -l ./ | tee /dev/stderr)"

lint:
	golangci-lint run

clean:
	$(GOCMD) clean

check: vet staticcheck gofmt

ci: build test check lint

bench:
	$(GOTEST) -bench=^Bench -count=5 -benchmem ./...
