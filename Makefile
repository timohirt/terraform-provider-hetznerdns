TEST?=$$(go list ./...)
GOFMT_FILES?=$$(find . -name '*.go')

BINARY_DIR=bin
BINARY_NAME=terraform-provider-hetznerdns

.PHONY: build testacc test lint fmt

build:
	mkdir -p $(BINARY_DIR)
	go build -o $(BINARY_DIR)/$(BINARY_NAME)

testacc:
	TF_LOG_PROVIDER=DEBUG TF_LOG=DEBUG TF_ACC=1 go test $(TEST) -v -timeout 180s

test: 
	go test $(TEST) || exit 1

lint:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	gofmt -w $(GOFMT_FILES)
