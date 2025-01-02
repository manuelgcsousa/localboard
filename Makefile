BINARY_NAME=localboard

all: build

build:
	go build -o bin/$(BINARY_NAME) cmd/$(BINARY_NAME)/main.go

run: build
	./bin/$(BINARY_NAME)

clean:
	rm -f bin/$(BINARY_NAME)
