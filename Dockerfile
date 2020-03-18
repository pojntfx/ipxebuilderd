# syntax=docker/dockerfile:experimental
# Build container
FROM --platform=$TARGETPLATFORM golang:alpine AS build
ARG DIBS_TARGET
ARG TARGETPLATFORM

WORKDIR /app

RUN apk add -u curl protoc git

RUN curl -Lo /tmp/dibs https://github.com/pojntfx/dibs/releases/latest/download/dibs-linux-amd64
RUN install /tmp/dibs /usr/local/bin

ENV GO111MODULE=on

RUN go get github.com/golang/protobuf/protoc-gen-go
RUN go get github.com/rakyll/statik
RUN go get github.com/mholt/archiver/cmd/arc

ADD . .

RUN dibs -generateSources
RUN dibs -build

# Run container
FROM --platform=$TARGETPLATFORM alpine
ARG DIBS_TARGET
ARG TARGETPLATFORM

RUN apk add -u gcc make git xz-dev perl musl-dev

COPY --from=build /app/.bin/binaries/ipxebuilderd* /usr/local/bin/ipxebuilderd

CMD /usr/local/bin/ipxebuilderd
