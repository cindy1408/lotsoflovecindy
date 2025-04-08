.PHONY: fmt
fmt:
	cd backend && gci write .
	cd backend && golangci-lint run --fix -v

.PHONY: start
start:
	cd backend && go run main.go & \
	cd ui && npm start
