# Aria2 Service

This is the Aria2 service

Generated with

```
micro new aria2-web --namespace=go.micro --alias=aria2 --type=web
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: go.micro.web.aria2
- Type: web
- Alias: aria2

## Dependencies

Micro services depend on service discovery. The default is multicast DNS, a zeroconf system.

In the event you need a resilient multi-host setup we recommend etcd.

```
# install etcd
brew install etcd

# run etcd
etcd
```

## Usage

A Makefile is included for convenience

Build the binary

```
make build
```

Run the service
```
./aria2-web
```

Build a docker image
```
make docker
```