# Build container
FROM golang:1.13.8 AS build

RUN apt update
RUN apt install -y protobuf-compiler

WORKDIR /app

ENV GO111MODULE=on

RUN go install github.com/golang/protobuf/protoc-gen-go
RUN go install github.com/rakyll/statik
RUN go install github.com/mholt/archiver/cmd/arc

COPY . .

RUN go generate ./...
RUN go build -o ipxebuilderd main.go

# Runner container
FROM debian:buster

RUN apt update
RUN apt install -y gcc make git liblzma-dev

COPY --from=build /app/ipxebuilderd /bin

CMD /bin/ipxebuilderd
