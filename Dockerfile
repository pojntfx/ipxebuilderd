# syntax=docker/dockerfile:experimental
FROM golang:1.13.8 AS build
WORKDIR /app
ARG TARGETPLATFORM

RUN apt update
RUN apt install -y protobuf-compiler

ENV GO111MODULE=on

RUN go install github.com/golang/protobuf/protoc-gen-go
RUN go install github.com/rakyll/statik
RUN go install github.com/mholt/archiver/cmd/arc
RUN curl -Lo dibs https://github.com/pojntfx/dibs/releases/latest/download/dibs-linux-amd64
RUN chmod +x dibs
RUN mv dibs /usr/local/bin/dibs

COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./.dibs.yml ./.dibs.yml
COPY ./main.go ./main.go
COPY ./cmd ./cmd
COPY ./pkg ./pkg

RUN dibs pipeline build assets

# Runner container
FROM --platform=$TARGETPLATFORM debian:buster-slim
ARG TARGETPLATFORM

RUN apt update
RUN apt install -y gcc make git liblzma-dev

COPY --from=build /app/.bin/ipxebuilderd-* /usr/local/bin/ipxebuilderd

EXPOSE 1440

CMD /bin/ipxebuilderd
