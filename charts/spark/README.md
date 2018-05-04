# Apache Spark Helm Chart

Apache Spark is a fast and general-purpose cluster computing system.

* http://spark.apache.org/

## Chart Details
This chart will do the following:

* 1 x Spark Master with port 8080 exposed
* 3 x Spark Workers with HorizontalPodAutoscaler to scale to max 10 pods when CPU hits 50% of 100m
* All using Kubernetes Deployments

## Prerequisites

* Assumes that serviceAccount tokens are available under hostname metadata. (Works on GKE by default) URL -- http://metadata/computeMetadata/v1/instance/service-accounts/default/token

## Installing the Chart

Download spark-0.1.0.tgz

To install the chart with the release name `my-release`:

```bash
$ helm install --name my-release spark-0.1.0.tgz
```

The command deploys spark on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

## Configuration

The following table lists the configurable parameters of the Spark chart and their default values.

### Spark Master

| Parameter               | Description                        | Default                                                    |
| ----------------------- | ---------------------------------- | ---------------------------------------------------------- |
| `Master.Name`           | Spark master name                  | `spark-master`                                             |
| `Master.Image`          | Container image name               | `p7hb/docker-spark`                                         |
| `Master.ImageTag`       | Container image tag                | `2.2.0`                                                 |
| `Master.Replicas`       | k8s deployment replicas            | `1`                                                        |
| `Master.Component`      | k8s selector key                   | `spark-master`                                             |
| `Master.Cpu`            | container requested cpu            | `100m`                                                     |
| `Master.Memory`         | container requested memory         | `512Mi`                                                    |
| `Master.ServicePort`    | k8s service port                   | `7077`                                                     |
| `Master.ContainerPort`  | Container listening port           | `7077`                                                     |
| `Master.DaemonMemory`   | Master JVM Xms and Xmx option      | `1g`                                                       |
| `Master.ServiceType `   | Kubernetes Service type            | `ClusterIP`                                             |

### Spark WebUi

|       Parameter       |           Description            |                         Default                          |
|-----------------------|----------------------------------|----------------------------------------------------------|
| `WebUi.Name`          | Spark webui name                 | `spark-webui`                                            |
| `WebUi.ServicePort`   | k8s service port                 | `8080`                                                   |
| `WebUi.ContainerPort` | Container listening port         | `8080`                                                   |

### Spark Worker

| Parameter                    | Description                        | Default                                                    |
| -----------------------      | ---------------------------------- | ---------------------------------------------------------- |
| `Worker.Name`                | Spark worker name                  | `spark-worker`                                             |
| `Worker.Image`               | Container image name               | `p7hb/docker-spark`                                         |
| `Worker.ImageTag`            | Container image tag                | `2.2.0`                                                 |
| `Worker.Replicas`            | k8s hpa and deployment replicas    | `3`                                                        |
| `Worker.ReplicasMax`         | k8s hpa max replicas               | `10`                                                       |
| `Worker.Component`           | k8s selector key                   | `spark-worker`                                             |
| `Worker.Cpu`                 | container requested cpu            | `100m`                                                     |
| `Worker.Memory`              | container requested memory         | `512Mi`                                                    |
| `Worker.ContainerPort`       | Container listening port           | `7077`                                                     |
| `Worker.CpuTargetPercentage` | k8s hpa cpu targetPercentage       | `50`                                                       |
| `Worker.DaemonMemory`        | Worker JVM Xms and Xmx setting     | `1g`                                                       |
| `Worker.ExecutorMemory`      | Worker memory available for executor | `1g`                                                       |

Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`.

Alternatively, a YAML file that specifies the values for the parameters can be provided while installing the chart. For example,

```bash
$ helm install --name my-release -f values.yaml stable/spark
```

> **Tip**: You can use the default [values.yaml](values.yaml)
