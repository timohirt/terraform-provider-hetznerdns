GO=/usr/local/go/bin/go
GOBUILD=$(GO) build
TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

BINARY_DIR=bin
BINARY_NAME=terraform-provider-hetznerdns

build:
	mkdir -p $(BINARY_DIR)
	$(GOBUILD) -o $(BINARY_DIR)/$(BINARY_NAME)

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v -timeout 30m

test: fmtcheck
	go test $(TEST) || exit 1

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	gofmt -w $(GOFMT_FILES)