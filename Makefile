.PHONY: build-app run clean test

build-app:
	go build -o bin/syncbuddy ./cmd/syncbuddy/main.go

run: build-app
	@./bin/syncbuddy

clean: 
	rm -rf bin/syncbuddy

test:
	go test ./... -cover