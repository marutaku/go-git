FROM golang:1.22


RUN apt-get update && apt-get install -y make bsdmainutils

WORKDIR /go/go-git/
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY go.mod .
COPY Makefile .

RUN go mod tidy && make all

CMD [ "/bin/bash" ]