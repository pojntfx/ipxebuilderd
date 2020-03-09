# ipxebuilderd

Build daemon and CLI for iPXE.

[![pipeline status](https://gitlab.com/pojntfx/ipxebuilderd/badges/master/pipeline.svg)](https://gitlab.com/pojntfx/ipxebuilderd/commits/master)

## Overview

`ipxebuilderd` is a build daemon with a CLI for [iPXE](https://ipxe.org), the leading open source network boot firmware. It is built of two components:

- `ipxebuilderd`, a build daemon with a gRPC interface which uploads built artifacts to S3
- `ipxectl`, a CLI for `ipxebuilderd`

`ipxebuilderd` bundles the iPXE source code, but does not the toolchain in it's binary. Please use to the Docker images or Helm chart, which bundle it; the latter also bundles a S3 ([Minio](https://min.io)) server.

## Installation

### Prebuilt Binaries

Prebuilt binaries are available on the [releases page](https://github.com/pojntfx/ipxebuilderd/releases/latest).

### Go Package

A Go package [is available](https://pkg.go.dev/github.com/pojntfx/ipxebuilderd).

### Docker Image

A Docker image is available on [Docker Hub](https://hub.docker.com/r/pojntfx/ipxebuilderd).

### Helm Chart

A Helm chart is available in [@pojntfx's Helm chart repository](https://pojntfx.github.io/charts/).

## Usage

### Daemon

You may also set the flags by setting env variables in the format `IPXEBUILDERD_[FLAG]` (i.e. `IPXEBUILDERD_IPXEBUILDERD_CONFIGFILE=examples/ipxebuilderd.yaml`) or by using a [configuration file](examples/ipxebuilderd.yaml).

```bash
% ipxebuilderd --help
ipxebuilderd is the iPXE build daemon.

Find more information at:
https://pojntfx.github.io/ipxebuilderd/

Usage:
  ipxebuilderd [flags]

Flags:
  -h, --help                                   help for ipxebuilderd
  -f, --ipxebuilderd.configFile string         Configuration file to use.
  -l, --ipxebuilderd.listenHostPort string     TCP listen host:port. (default "0.0.0.0:1440")
  -u, --ipxebuilderd.s3AccessKey string        Access key of the S3 server to connect to. (default "ipxebuilderUser")
  -b, --ipxebuilderd.s3Bucket string           S3 bucket to use. (default "ipxebuilderd")
  -s, --ipxebuilderd.s3HostPort string         Host:port of the S3 server to connect to. (default "minio.ipxebuilderd.felix.pojtinger.com")
  -o, --ipxebuilderd.s3HostPortPublic string   Public host:port of the S3 server (will be used in shared links). (default "minio.ipxebuilderd.felix.pojtinger.com")
  -p, --ipxebuilderd.s3SecretKey string        Secret key of the S3 server to connect to. (default "ipxebuilderdPass")
  -z, --ipxebuilderd.secure                    Whether to use a secure connection to S3.
```

### Client CLI

You may also set the flags by setting env variables in the format `IPXE_[FLAG]` (i.e. `IPXE_IPXE_CONFIGFILE=examples/ipxe.yaml`) or by using a [configuration file](examples/ipxe.yaml).

```bash
% ipxectl --help
ipxectl manages ipxebuilderd, the iPXE build daemon.

Find more information at:
https://pojntfx.github.io/ipxe/

Usage:
  ipxectl [command]

Available Commands:
  apply       Apply a ipxe
  delete      Delete one or more iPXE(s)
  get         Get one or all iPXE(s)
  help        Help about any command

Flags:
  -h, --help   help for ipxectl

Use "ipxectl [command] --help" for more information about a command.
```

## License

ipxebuilderd (c) 2020 Felix Pojtinger

SPDX-License-Identifier: AGPL-3.0
