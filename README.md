# Signa

Signa is a Slack bot and ChatOps kit written in Go.

The bot offers a framework to develop ChatOps tools and it currently has
two extensions for Kubernetes.

## Install

```
$ go get [-u] github.com/signavio/signa/cmd/signa
```

## Extensions

Below you may check a list of the current available extensions.

### kubernetes/deployment

[ext/kubernetes/deployment](ext/kubernetes/deployment) is a deployment
command that runs on top of `kubectl`.

### kubernetes/get

[ext/kubernetes/get](ext/kubernetes/get) retrieves information of the
resources running in the cluster.

## Internals

In this section you may find an overview of the internals of Signa and
code organization.

Signa has three main components:
- [cmd](cmd/) is the home of the executable code. It's the entrypoint to install Signa.
- [ext](ext/) is where all the native extensions are placed.
- [pkg](pkg/) is the Signa core libraries and utils.

## TODO

- Add test suite with different test cases.
- Add in-depth documentation.

## Maintainers

Stephano Zanzin - [@microwaves](https://github.com/microwaves)
