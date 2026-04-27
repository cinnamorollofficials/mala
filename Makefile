run:
	go run main.go

build:
	go build -o bin/mala main.go

tidy:
	go mod tidy

test:
	go test ./...

clean:
	rm -rf bin/

perf-test:
	go run scripts/perf_test.go

# Docker commands
docker-up:
	docker compose up -d

docker-down:
	docker compose down

docker-logs:
	docker compose logs -f

docker-build:
	docker compose build