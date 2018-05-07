# Kubeflow

The [Kubeflow](https://github.com/kubeflow/kubeflow) project is dedicated to making deployments of machine learning (ML) workflows on Kubernetes simple, portable and scalable. Our goal is not to recreate other services, but to provide a straightforward way to deploy best-of-breed open-source systems for ML to diverse infrastructures. Anywhere you are running Kubernetes, you should be able to run Kubeflow.

## Introduction

This chart is used for the Kubeflow benchmark on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.9+

## Installing the Chart

Download kubeflow-0.1.0.tgz

To install the chart with the release name `my-release`:

```bash
$ helm install --name my-release kubeflow-0.1.0.tgz
```

The command deploys Kubeflow on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `my-release` deployment:

```bash
$ helm delete my-release
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Kubeflow chart and their default values.

### Ambassador

| Parameter               | Description                        | Default                                                    |
| ----------------------- | ---------------------------------- | ---------------------------------------------------------- |
| `ambassador.limitsCpu`           | CPU resource limits                  | `1`                                             |
| `ambassador.limitsMemory`          | Memory resource limits               | `400Mi`                                         |
| `ambassador.requestsCpu`       | CPU resource requests                | `200m`                                                 |
| `ambassador.requestsMemory`       | Memory resource requests            | `100Mi`                                                        |
| `ambassador.image`      | ambassador image                   | `quay.io/datawire/ambassador:0.30.1`                                             |
| `ambassador.statsdImage`            | ambassador statsdImage            | `quay.io/datawire/statsd:0.30.1`                                                     |
| `ambassador.replicaCount`         | number of ambassador replica         | `3`                                                    |
| `ambassador.adminServicePort`    | port of ambassador admin service                   | `8877`                                                     |
| `ambassador.adminServiceType`  | type of ambassador admin service           | `ClusterIP`                                                     |
| `ambassador.servicePort`   | port of ambassador service      | `80`                                                       |
| `ambassador.serviceType`   | type of ambassador service            | `ClusterIP`                                             |

### Spartakus

|       Parameter       |           Description            |                         Default                          |
|-----------------------|----------------------------------|----------------------------------------------------------|
| `spartakus.replicaCount`          | number of spartakus replica                 | `1`                                            |
| `spartakus.image`   | spartakus image                 | `gcr.io/google_containers/spartakus-amd64:v1.0.0`                                                   |

### TFJob

| Parameter                    | Description                        | Default                                                    |
| -----------------------      | ---------------------------------- | ---------------------------------------------------------- |
| `tfjob.replicaCount`                | number of tfjob replica                  | `1`                                             |
| `tfjob.image`               | tfjob image               | `gcr.io/kubeflow-images-public/tf_operator:v20180329-a7511ff`                                         |

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
$ helm install --name my-release -f values.yaml kubeflow-0.1.0.tgz
```

> **Tip**: You can use the default [values.yaml](values.yaml)