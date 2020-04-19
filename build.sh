#!/usr/bin/env bash

function proto() {
  arr=(aria2-srv movie-srv)
  for d in "${arr[@]}"; do
    for f in "$d"/proto/**/*.proto; do
      protoc --proto_path="${GOPATH}"/src:. --micro_out=. --go_out=. "$f"
      echo compiled: "$f"
    done
  done
}

function run() {
  pushd "$1" >/dev/null
  go run main.go plugin.go
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

ENV="MICRO_REGISTRY=etcd MICRO_REGISTRY_ADDRESS=192.168.1.105:2379"

case $1 in
proto)
  proto
  ;;
run)
  run "$2"
  ;;
api-web)
  MICRO_API_HANDLER=web \
  eval "${ENV}" micro --cors-allowed-headers="Origin,Content-Type,Accept,Authorization"  --cors-allowed-origins="localhost:8081"  --cors-allowed-methods="HEAD,GET,POST,OPTIONS,PUT" api --enable_cors true
  ;;
api-api)
  MICRO_API_HANDLER=api \
  eval "${ENV}" micro api
  ;;
web)
  eval "${ENV}" micro web
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
