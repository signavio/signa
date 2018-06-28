# Signa

Signa is a Slack bot and ChatOps tool written in Go.

It offers built-in commands for Kubernetes and a framework to develop ChatOps tools.

## Install

```
$ go get [-u] github.com/signavio/signa/cmd/signa
```

## Configuration

Refer to the [signa.sample.yaml](signa.sample.yaml) configuration file to setup Signa.

## Extensions

Find a list of the current built-in commands below.

### kubernetes/deployment

[kubernetes/deployment](ext/kubernetes/deployment) is a deployment
command that runs on top of `kubectl`.

Usage example:

```
!deploy app-backend app-backend-container cluster1 v1.0.1
```

Where `app-backend` stands for the application name, `app-backend-container` for the container name,
`cluster1` for the cluster name and `v1.0.1` for the desired image tag.

### kubernetes/get

[kubernetes/get](ext/kubernetes/get) retrieves information of the
resources running in the cluster.

Usage example:

```
!get pods -n foobar-namespace
```

Currently it only works for 1 single context, the support to multi-contexts/clusters should be added.

### kubernetes/jobs

[kubernetes/jobs](ext/kubernetes/jobs) runs kubernetes jobs using a pre-determined configuration.

Usage example:

```
!run make-pizza cluster1 v0.9.7
```

Where `make-pizza` stands for the job name, `cluster1` for the cluster name and `v0.9.7` for the
image tag.

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
