# Capstan deploying

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Setup Kubernetes cluster](#kubernetes-cluster)
    - [GKE](#google-kubernetes-engine)
- [Install packages](#install-packages)
    - [Install Docker](#install-docker)
    - [Install Pushgateway](#install-pushgateway)
    - [Install Prometheus](#install-prometheus)
    - [Install Grafana](#install-grafana)
    - [Install capstan](#install-capstan)

## Overview

This document shows how to start the capstan benchmarker with a kubernetes cluster offered by a cloud provider.

## Prerequisites

Before you can run the capstan Benchmarker, you need an account on the cloud provider you want to benchmark it's Kubernetes service:

- [GCP](https://cloud.google.com)
- [AWS](http://aws.amazon.com)
- [Azure](http://azure.microsoft.com)
- [Aliyun](http://www.aliyun.com)

You also need the software dependencies, which are mostly command line tools and credentials to access your accounts without a password. The following steps will help you with this on a cloud provider.

## Kubernetes cluster

This section helps you to setup a Kubernetes cluster by the cloud provider's Kubernetes service.

- [GKE](#google-kubernetes-engine)

### Google Kubernetes Engine

1. Install `gcloud`

    Follow the [instructions](https://developers.google.com/cloud/sdk/) and install the `gcloud`.

1. setup authentication:

    ```bash
    gcloud init
    ```

1. setup kubernetes cluster:

    ```bash
    gcloud container clusters create example-cluster
    ```

For more gcloud help, see [`gcloud` docs](https://cloud.google.com/sdk/gcloud/).

## Install packages

Docker, Pushgateway, Prometheus, Grafana and capstan should be installed on the same machine.

### Install Docker

On Ubuntu 16.04+:

```sh
apt-get update
apt-get install -y docker.io
```

On CentOS 7:

```sh
yum install -y docker
```

Configure and start docker:

```sh
systemctl enable docker
systemctl start docker
```

### Install Pushgateway

```sh
docker run -d -p 9091:9091 prom/pushgateway
```

### Install Node-exporter

```sh
Install Node-exporter inside your kubernetes.

cat >/tmp/node-exporter.yaml <<EOF
apiVersion: extensions/v1beta1
kind: DaemonSet
metadata:
  name: node-exporter
  labels:
    k8s-app: node-exporter
  namespace: exporter
spec:
  template:
    metadata:
      labels:
        k8s-app: node-exporter
    spec:
      containers:
      - image: prom/node-exporter
        name: node-exporter
        ports:
        - containerPort: 9100
          protocol: TCP
          name: http
---
apiVersion: v1
kind: Service
metadata:
  labels:
    k8s-app: node-exporter
  name: node-exporter
  namespace: exporter
spec:
  ports:
  - name: http
    port: 9100
    nodePort: 31672
    protocol: TCP
  type: NodePort
  selector:
    k8s-app: node-exporter
EOF

kubectl create namespace exporter

kubectl create -f /tmp/node-exporter.yaml
```

### Install Prometheus

Configure Prometheus:

```sh
cat >/tmp/prometheus.yml <<EOF
global:
  scrape_interval: 15s
  scrape_timeout: 10s
scrape_configs:
  - job_name: 'pushgateway'
    static_configs:
      - targets: ['<Your-HostIP>:9091']
  - job_name: 'node-exporter'
    static_configs:
      - targets: ['<Your-Kubernetes-MasterIP>:31672','<Your-Kubernetes-Node1IP>:31672','<Your-Kubernetes-Node2IP>:31672',...]
EOF
```

Start Prometheus:

```sh
docker run -d -p 9090:9090 -v /tmp/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```

### Install Grafana

```sh
mkdir -p /tmp/provisioning/datasources

cat >/tmp/provisioning/datasources/prometheus.yaml <<EOF
apiVersion: 1
datasources:
  - name: prometheus
    type: prometheus
    access: proxy
    url: http://<Your-HostIP>:9090
    basicAuth: false
    editable: true
EOF

mkdir /tmp/provisioning/dashboards

cp -r $GOPATH/src/github.com/ZJU-SEL/capstan/grafana-dashboards/* /tmp/provisioning/dashboards/

docker run -d -p 3000:3000 -v /tmp/provisioning:/etc/grafana/provisioning grafana/grafana
```

### Install capstan

Build capstan:

```sh
mkdir -p $GOPATH/src/github.com/ZJU-SEL
git clone https://github.com/ZJU-SEL/capstan.git $GOPATH/src/github.com/ZJU-SEL/capstan
cd $GOPATH/src/github.com/ZJU-SEL/capstan
make && make install
```

Configure capstan:

```sh
cat >/etc/capstan/config <<EOF
{
    "ResultsDir": "/tmp/capstan",
    "Provider": "aliyun",
    "Prometheus": {
        "PushgatewayEndpoint": "http://<Your-HostIP>:9091"
    },
    "Steps": 10,
    "Namespace": "capstan",
    "Workloads": [
        {
            "name": "nginx",
            "image": "nginx:1.7.9",
            "frequency": 5,
            "testingTool": {
                "name": "wrk",
                "image": "wadelee/wrk",
                "steps": 10,
                "testingCaseSet": [
                    {
                        "name": "benchmarkPodIPDiffNode",
                        "testingToolArgs": "-t10 -c100 -d90 http://\$(ENDPOINT)/"
                    },
                    {
                        "name": "benchmarkPodIPSameNode",
                        "testingToolArgs": "-t10 -c100 -d90 http://\$(ENDPOINT)/"
                    }
                ]
            }
        }
    ]
}
EOF
```

Start capstan:

```sh
capstan --v=3 --logtostderr --config=/etc/capstan/config --kubeconfig=/etc/kubernetes/admin.conf &
```