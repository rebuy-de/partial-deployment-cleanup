# partial-deployment-cleanup

Purges old branches and deployments

## Requirements

* `glide` https://github.com/Masterminds/glide#install
* Go >= `1.5`
  + if you have version `1.5`, then you have to `export GO15VENDOREXPERIMENT=1`

## Build

1. download dependencies: `glide install`
2. build binary: `go build`
3. use fancy binary: `./partial-deployment-cleanup`

## Update dependencies

1. update dependencies: `glide update`
2. test everything: `go test ./...`
3. debug: `vim`

## Use

* `./partial-deployment-cleanup -h`
