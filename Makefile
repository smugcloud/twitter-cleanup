.PHONY: build-linux, build-darwin

all:
	$(MAKE) build-linux
	$(MAKE) build-darwin 

build-linux:
	GOOS=linux GOARCH=amd64 go build -o builds/linux/twitter-cleanup

build-darwin:
	GOOS=darwin GOARCH=amd64 go build -o builds/darwin/twitter-cleanup


test:
	docker build -f Dockerfile -t twitter-cleanup:testing .
	docker run --rm --name twitter-cleanup twitter-cleanup:testing
	$(MAKE) build-linux