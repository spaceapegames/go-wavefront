TEST?=$$(go list ./... |grep -v 'vendor'|grep -v 'examples')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: test

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	gofmt -w $(GOFMT_FILES)

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
		echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'" | grep -v example

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test -v $(TESTARGS) -timeout=30s -parallel=4
