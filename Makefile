.PHONY: fmt
fmt:
	@echo "→ Running gci and golangci-lint in ./backend"
	cd api && gci write .
	cd api && golangci-lint run --fix -v ./...

.PHONY: start
start:
	@echo "→ Starting Go backend and UI"
	cd api && go run main.go & \
	cd ui && npm start

.PHONY: rebuild
rebuild:
	docker-compose down && docker-compose build --no-cache && docker-compose up

