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