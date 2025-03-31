.PHONY: fmt
fmt:
	gci write .
	golangci-lint run --fix -v
