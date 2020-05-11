CC=go
VERSION=$(shell cat VERSION)
OS=$(shell uname -s)
ARCH=$(shell uname -m)
BINARY=eveleve
include config.make

linux: bot-linux 
darwin: bot-darwin 
windows: bot-windows 

bot-linux: bin/linux/$(BINARY)
bot-windows: bin/windows/$(BINARY).exe
bot-darwin: bin/darwin/$(BINARY)

all: linux windows darwin test

SOURCES=main.go \
		config.go \
		master.go \
		github.go \
		travis.go \
		project.go \
		patreon.go \
		discord.go \
		status.go

test:
	$(CC) test . -v

packet.pb.go: packet.proto
	protoc --go_out=. $^

bin/linux/$(BINARY): $(SOURCES)
	GOOS=linux $(CC) build -ldflags="-w -s -X main.AppVersion=$(APP_VERSION)" -o $@ -v $^

bin/windows/$(BINARY).exe: $(SOURCES)
	GOOS=windows $(CC) build -ldflags="-w -s -X main.AppVersion=$(APP_VERSION)" -o $@ -v $^
	
bin/darwin/$(BINARY): $(SOURCES)
	GOOS=darwin $(CC) build -ldflags="-w -s -X main.AppVersion=$(APP_VERSION)" -o $@ -v $^


clean:
	-rm -f bin/linux/$(BINARY)
	-rm -f bin/windows/$(BINARY).exe
	-rm -f bin/darwin/$(BINARY)

mrproper:
	-rm -rf bin
	-rm -f config.make
	-rm -f packet.pb.go
