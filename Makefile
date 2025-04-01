.PHONY: fmt
fmt:
	gci write .
	golangci-lint run --fix -v

.PHONY: run
run:
	cd backend && go run main.go & \
	cd ui && npm start