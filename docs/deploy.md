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
    "UUID": "123456",
    "ResultsDir": "/tmp/capstan",
    "Provider": "aliyun",
    "Address": "0.0.0.0:8080",
    "PushgatewayEndpoint": "http://<Your-HostIP>:9091",
    "Steps": 10,
    "Namespace": "capstan",
    "Workloads": [
        {
            "name": "nginx",
            "helm": {
                "name": "chart1",
                "set": "imageTag=1.7.9",
                "chart": "charts/nginx-0.1.0.tgz"
            },
            "frequency": 5,
            "testTool": {
                "name": "wrk",
                "script": "scripts/wrk.sh",
                "steps": 10,
                "testCaseSet": [
                    {
                        "name": "case1",
                        "affinity": true,
                        "args": "-t10 -c100 -d30 http://chart1-nginx/",
                        "metrics": "QPS"
                    },
                    {
                        "name": "case2",
                        "affinity": false,
                        "args": "-t10 -c100 -d30 http://chart1-nginx/",
                        "metrics": "QPS"
                    }
                ]
            }
        },
        {
            "name": "iperf3",
            "helm": {
                "name": "chart2",
                "chart": "charts/iperf3-0.1.0.tgz"
            },
            "frequency": 5,
            "testTool": {
                "name": "iperf3",
                "script": "scripts/iperf3.sh",
                "steps": 10,
                "testCaseSet": [
                    {
                        "name": "case1",
                        "affinity": true,
                        "args": "-c chart2-iperf3",
                        "metrics": "BandWidth"
                    },
                    {
                        "name": "case2",
                        "affinity": false,
                        "args": "-c chart2-iperf3",
                        "metrics": "BandWidth"
                    }
                ]
            }
        },
        {
            "name": "mysql",
            "helm": {
                "name": "chart3",
                "set": "mysqlRootPassword=capstan,persistence.enabled=false",
                "chart": "stable/mysql"
            },
            "frequency": 5,
            "testTool": {
                "name": "tpcc-mysql",
                "script": "scripts/tpcc-mysql.sh",
                "steps": 10,
                "testCaseSet": [
                    {
                        "name": "case1",
                        "affinity": true,
                        "args": "-w1 -c10 -r60 -l60",
                        "envs": "MYSQL_HOST=chart3-mysql,MYSQL_ROOT_PASSWORD=capstan",
                        "metrics": "TPMC"
                    },
                    {
                        "name": "case2",
                        "affinity": false,
                        "args": "-w1 -c10 -r60 -l60",
                        "envs": "MYSQL_HOST=chart3-mysql,MYSQL_ROOT_PASSWORD=capstan",
                        "metrics": "TPMC"
                    }
                ]
            }
        }
    ]
}
EOF
```

### Install Pushgateway

```sh
docker run -d -p 9091:9091 prom/pushgateway
```

### Install Prometheus

Configure Prometheus:

Now we run the Kubernetes on Aliyun. So here we use Aliyun as example. It is sure that you can also use GCP, AWS, and Azure.

```sh
cat >/etc/capstan/prometheus/prometheus.yml <<EOF
global:
  scrape_interval: 15s
  scrape_timeout: 10s
scrape_configs:
  - job_name: 'pushgateway'
    static_configs:
      - targets: ['<Your-HostIP>:9091']
  - job_name: 'aliyun'
    static_configs:
      - targets: ['<Your-Kubernetes in aliyun-MasterIP>:31672','<Your-Kubernetes in aliyun-Node1IP>:31672','<Your-Kubernetes in aliyun-Node2IP>:31672',...]
EOF
```

Start node-exporter:

```sh
kubectl create namespace capstan-exporter
kubectl create -f $GOPATH/src/github.com/ZJU-SEL/capstan/deploy/node-exporter.yaml
```

Start Prometheus:

```sh
docker run -d -p 9090:9090 -v /etc/capstan/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
```

### Install Grafana

```sh

IPADDR=$(curl -s http://members.3322.org/dyndns/getip)
sed -i "s/<Your-HostIP>/${IPADDR}/g" /etc/capstan/grafana/provisioning/datasources/prometheus.yaml /etc/capstan/config

sh $GOPATH/src/github.com/ZJU-SEL/capstan/grafana-dashboards/configDashboard.sh

docker run -d -p 3000:3000 -v /etc/capstan/grafana/provisioning:/etc/grafana/provisioning grafana/grafana
```

### Start capstan

```sh
capstan --v=3 --logtostderr --config=/etc/capstan/config --kubeconfig=/etc/kubernetes/admin.conf &
```