BINARY_NAME=localboard
SOURCE_FILES=*.go

all: build

build:
	go build -o bin/$(BINARY_NAME) $(SOURCE_FILES)

run: build
	./bin/$(BINARY_NAME)

clean:
	rm -f bin/$(BINARY_NAME)
