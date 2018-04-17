# Nginx

[Nginx](https://nginx.org/en/)[engine x] is an HTTP and reverse proxy server, a mail proxy server, and a generic TCP/UDP proxy server, originally written by Igor Sysoev. For a long time, it has been running on many heavily loaded Russian sites including Yandex, Mail.Ru, VK, and Rambler.

## Introduction

This chart bootstraps a single node nginx deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.9+

## Installing the Chart

Download nginx-0.1.0.tgz

To install the chart with the release name `my-release`:

```bash
$ helm install --name my-release nginx-0.1.0.tgz
```

The command deploys nginx on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Nginx chart and their default values.

| Parameter                            | Description                               | Default                                              |
| ------------------------------------ | ----------------------------------------- | ---------------------------------------------------- |
| `imageTag`                           | `nginx` image tag.                        | Most recent release                                  |
| `imagePullPolicy`                    | Image pull policy                         | `IfNotPresent`                                       |
| `resources.limits.memory`                          | Memory resource limits       | `128Mi`                         |
| `resources.limits.cpu`                          | CPU resource limits       | `100m`                         |
| `resources.requests.memory`                          | Memory resource requests       | `128Mi`                         |
| `resources.requests.cpu`                          | CPU resource requests       | `100m`                         |

Some of the parameters above map to the env variables defined in the [Nginx DockerHub image](https://hub.docker.com/_/nginx/).

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
$ helm install --name my-release -f values.yaml nginx-0.1.0.tgz
```

> **Tip**: You can use the default [values.yaml](values.yaml)