# partial-deployment-cleanup

Purges old branches and deployments

## Primary responsible person (PRP)

* [@prp-partial-deployment-cleanup](https://github.com/orgs/rebuy-de/teams/prp-partial-deployment-cleanup)

## Requirements

* `go` >= `1.6` for building (https://golang.org/)
* `consul` for testing (https://www.consul.io/)
* `fpm` for packaging (https://github.com/jordansissel/fpm)

## Testing

```bash
hack/test.sh
```

## Build

```bash
hack/build.sh
```

## Build RPM

```bash
hack/package.sh
```

* This gets the version from git, so make sure you have a proper tag and no 
  uncommited changes in your workspace. Otherwise the version string will look
  ugly (but still unambiguous).

## Use

```bash
target/partial-deployment-cleanup -h
```
