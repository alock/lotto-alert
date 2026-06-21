BINARY := lotto-alert

.PHONY: build test pi lint clean

build:
	go build -o $(BINARY)

test:
	go test ./...

pi:
	GOOS=linux GOARCH=arm GOARM=7 go build -o $(BINARY)
	scp $(BINARY) config/email.json pi.local:/home/pi/.local/bin
	rm $(BINARY)

lint:
	go vet ./...

clean:
	rm -f $(BINARY)
