# Aria2 Service

This is the Aria2 service

Generated with

```
micro new aria2-srv --namespace=micro.srv.aria2 --type=service
```

## Getting Started

- [Configuration](#configuration)
- [Dependencies](#dependencies)
- [Usage](#usage)

## Configuration

- FQDN: micro.srv.aria2.service.aria2
- Type: service
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
./aria2-service
```

Build a docker image
```
make docker
```