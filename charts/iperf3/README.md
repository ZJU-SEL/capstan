# iPerf3

[iPerf3](https://iperf.fr/) is a tool for active measurements of the maximum achievable bandwidth on IP networks. It supports tuning of various parameters related to timing, buffers and protocols (TCP, UDP, SCTP with IPv4 and IPv6). For each test it reports the bandwidth, loss, and other parameters.

## Introduction

This chart bootstraps a single node iPerf3 server deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.9+

## Installing the Chart

Download iperf3-0.1.0.tgz

To install the chart with the release name `my-release`:

```bash
$ helm install --name my-release iperf3-0.1.0.tgz
```

The command deploys iPerf3 on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the iPerf3 chart and their default values.

| Parameter                            | Description                               | Default                                              |
| ------------------------------------ | ----------------------------------------- | ---------------------------------------------------- |
| `imageTag`                           | `iPerf3` image tag.                        | Most recent release                                  |
| `imagePullPolicy`                    | Image pull policy                         | `IfNotPresent`                                       |
| `resources.limits.memory`                          | Memory resource limits       | `128Mi`                         |
| `resources.limits.cpu`                          | CPU resource limits       | `200m`                         |
| `resources.requests.memory`                          | Memory resource requests       | `128Mi`                         |
| `resources.requests.cpu`                          | CPU resource requests       | `200m`                         |

Some of the parameters above map to the env variables defined in the [iPerf3 DockerHub image](https://hub.docker.com/r/wadelee/iperf3/).

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
$ helm install --name my-release -f values.yaml iperf3-0.1.0.tgz
```

> **Tip**: You can use the default [values.yaml](values.yaml)