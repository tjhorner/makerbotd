.PHONY: dist dist-win dist-macos dist-linux-amd64 dist-linux-arm ensure-dist-dir build install uninstall

GOBUILD=go build -ldflags="-s -w"
INSTALLPATH=/usr/local/bin

ensure-dist-dir:
	@- mkdir -p dist

dist-win: ensure-dist-dir
	# Build for Windows x64
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o dist/makerbotd-windows-amd64.exe *.go

dist-macos: ensure-dist-dir
	# Build for macOS x64
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o dist/makerbotd-darwin-amd64 *.go

dist-linux-amd64: ensure-dist-dir
	# Build for Linux x64
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o dist/makerbotd-linux-amd64 *.go

dist-linux-arm: ensure-dist-dir
	# Build for Linux ARM
	GOOS=linux GOARCH=arm $(GOBUILD) -o dist/makerbotd-linux-arm *.go

dist: dist-win dist-macos dist-linux-amd64 dist-linux-arm

build:
	@- mkdir -p bin
	$(GOBUILD) -o bin/makerbotd *.go
	@- chmod +x bin/makerbotd

install: build
	mv bin/makerbotd $(INSTALLPATH)/makerbotd
	@- rm -rf bin
	@echo "makerbotd was installed to $(INSTALLPATH)/makerbotd. Run make uninstall to get rid of it, or just remove the binary yourself."

install-systemd: build
	mv bin/makerbotd $(INSTALLPATH)/makerbotd
	@- rm -rf bin
	cp makerbotd.service /etc/systemd/system/makerbotd.service
	@echo "makerbotd was installed to $(INSTALLPATH)/makerbotd and the service file was installed to /etc/systemd/system/makerbotd.service. To start makerbotd, run `systemctl start makerbotd`. You can find the config file in /etc/makerbotd/config.json."

uninstall:
	rm $(INSTALLPATH)/makerbotd

run:
	@- go run *.go