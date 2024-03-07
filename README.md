<h1 style="color:#f9532f">Tiger Beetle gRPC</h1>

A simple gRPC server for TigerBeetle

<img width=200 src="https://tigerbeetle.com/60f5f501f4be1be5f45a50a3_img-performance.png">

Proto File: [/proto/tigerbeetle.proto](/proto/tigerbeetle.proto)

Config Example File: [/config-example.yml](/config-example.yml)

## Development setup

**1. Install [taskfile](https://taskfile.dev/installation/)**

**2. Run setup tasks**

```
$ task setup
```

Setup and run tigerbeetle in docker (optional)

```
$ task docker:setup docker:start
```

**3. Copy example config file**

```
$ cp config-example.yml config.yml
```

**4. Run server with the following command**

```
$ task start
```

Or run and watch for changes

```
$ task dev
```

## License

Apache License Version 2.0