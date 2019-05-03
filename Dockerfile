FROM golang:1.12

WORKDIR /twitter-cleanup

COPY . /twitter-cleanup

CMD ["/usr/local/go/bin/go", "test", "./..."]