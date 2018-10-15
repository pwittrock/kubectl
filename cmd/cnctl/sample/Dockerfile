

FROM golang:1.11.1

RUN mkdir -p $GOPATH/src/webserver
COPY main.go $GOPATH/src/webserver/main.go
WORKDIR $GOPATH/src/webserver
RUN go build main.go

CMD $GOPATH/src/webserver/main