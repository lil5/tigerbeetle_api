<h1 style="color:#f9532f">Tiger Beetle REST</h1>

A simple REST api server for [TigerBeetle](https://tigerbeetle.com/)

<img width=200 src="https://tigerbeetle.com/60f5f501f4be1be5f45a50a3_img-performance.png">

**Rest API:** Tigerbeetle REST API uses 🐶 Bruno

Download the client here: https://www.usebruno.com/
And get started by opening the `/bruno` directory in Bruno.

**Config Example File:** [/config-example.yml](/config-example.yml)

## Development setup

**1. Install [taskfile](https://taskfile.dev/installation/) and [golang](https://go.dev/)**

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