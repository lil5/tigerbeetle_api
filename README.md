[![GoDoc](https://pkg.go.dev/badge/github.com/lil5/tigerbeetle_api?status.svg)](https://pkg.go.dev/github.com/lil5/tigerbeetle_api?tab=doc)

<h1 style="color:#f9532f">Tiger Beetle REST</h1>

A simple REST api server for [TigerBeetle](https://tigerbeetle.com/)

<img width=200 src="/screenshot_bruno.webp">

**Rest API:** Tigerbeetle REST API uses 🐶 Bruno

Download the client here: https://www.usebruno.com/
And get started by opening the `/bruno` directory in Bruno.

**Config Example File:** [/config-example.yml](/config-example.yml)

## Development setup

**1. Install [taskfile](https://taskfile.dev/installation/) and [golang](https://go.dev/)**

**2. Install zig for cross-compilation**

```
$ brew install zig
```

Setup and run tigerbeetle in docker (optional)

```
$ make docker_setup docker_start
```

**3. Copy example config file**

```
$ cp .example.env .env
```

**4. Run server with the following command**

```
$ make start
```

## License

Apache License Version 2.0