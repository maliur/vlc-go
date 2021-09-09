run:
	go build -o ./vlcgo examples/status/main.go && ./vlcgo

playlist:
	go build -o ./vlcgo examples/playlist/main.go && ./vlcgo

test:
	go test

lint:
	golangci-lint run ./...

.PHONY: clean
clean:
	rm -f ./vlcgo
