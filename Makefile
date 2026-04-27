run:
	go run main.go

build:
	go build -o bin/mala main.go

tidy:
	go mod tidy

# migrate:
# 	go 

test:
	go test ./...

clean:
	rm -rf bin/

perf-test:
	go run scripts/perf_test.go