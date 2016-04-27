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
hack/test.sh
```

## Build

```bash
hack/build.sh
```

## Use

```bash
./partial-deployment-cleanup -h
```
