all: build_linux build_osx

build_linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o tunnel_linux

build_osx:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build  -o tunnel_osx

.PHONY: clean
clean:
	rm -f tunnel_linux  tunnel_osx
