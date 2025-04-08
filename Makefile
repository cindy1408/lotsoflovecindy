.PHONY: fmt
fmt:
	@echo "→ Running gci and golangci-lint in ./backend"
	cd backend && gci write .
	cd backend && golangci-lint run --fix -v ./...

.PHONY: start
start:
	@echo "→ Starting Go backend and UI"
	cd backend && go run main.go & \
	cd ui && npm start