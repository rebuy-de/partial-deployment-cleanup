# partial-deployment-cleanup

Purges old branches and deployments

## Primary responsible person (PRP)

* [@prp-partial-deployment-cleanup](https://github.com/orgs/rebuy-de/teams/prp-partial-deployment-cleanup)

## Requirements

* Go >= `1.6`
* `glide` https://github.com/Masterminds/glide#install
* since the test a starting a real Consul server, `consul` is required in your `$PATH`

## Testing

```bash
go test $(glide novendor) -p 1
```

* `glide novendor` lists all subpackages; otherwise `go` would also test `vendor/`
* you have to disable parallel testing with `-p 1`; this is a limitation of `github.com/hashicorp/consul/testutil`

## Build

1. download dependencies: `glide install`
2. build binary: `go build`
3. use fancy binary: `./partial-deployment-cleanup`

## Update dependencies

1. update dependencies: `glide update`
2. test everything: `go test $(glide novendor)`
3. debug: `vim`

## Use

```bash
./partial-deployment-cleanup -h
```
