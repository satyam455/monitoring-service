.PHONY: fmt test run compose-up compose-down

fmt:
	gofmt -w cmd internal

test:
	go test ./...

run:
	go run ./cmd/monitor

compose-up:
	docker compose up --build

compose-down:
	docker compose down
