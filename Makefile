Version := $(shell git describe --tags --abbrev=0)

.PHONY: dist
dist:
	mkdir -p bin/
	GOOS=windows GOARCH=amd64 go build -o bin/zimg-${Version}.exe
	GOOS=darwin GOARCH=arm64 go build -o bin/zimg-${Version}

.PHONY: build
build:
	go build -o bin/zimg-test
	mv bin/zimg-test $(HOME)/bin/zimg
