#!/usr/bin/env bash

function proto() {
  arr=(aria2-srv movie-srv)
  GOPATH=/c/Users/admin/go/
  for d in "${arr[@]}"; do
    for f in "$d"/proto/**/*.proto; do
      protoc --proto_path="${GOPATH}"/src:. --micro_out=. --go_out=. "$f"
      echo compiled: "$f"
    done
  done
}

function run() {
  pushd "$1" >/dev/null
  micro run service
  popd >/dev/null
}

function build() {
  arr=(config-grpc-srv aria2-srv movie-srv)
  for d in "${arr[@]}"; do
    pushd "$d" >/dev/null
    GOOS=linux GOARCH=$1 CGO_ENABLED=0 go build -o "$d" -ldflags "-w -s" main.go plugin.go
    echo compiled: "$d"
    popd >/dev/null
  done
}

case $1 in
proto)
  proto
  ;;
run)
  run "$2"
  ;;
api-http)
  MICRO_REGISTRY=etcd \
  MICRO_REGISTRY_ADDRESS=192.168.1.105:2379 \
  MICRO_API_HANDLER=http \
  MICRO_API_NAMESPACE=go.micro.web \
  micro api
  ;;
api-api)
  MICRO_REGISTRY=etcd \
  MICRO_REGISTRY_ADDRESS=192.168.1.105:2379 \
  MICRO_API_HANDLER=api \
  MICRO_API_NAMESPACE=go.micro.api \
  micro api
  ;;
web)
  MICRO_REGISTRY=etcd \
  MICRO_REGISTRY_ADDRESS=192.168.1.105:2379 \
  micro web
  ;;
build)
  build amd64
  ;;
build-arm64)
  build arm64
  ;;
*)
  echo please check you input
  ;;
esac
