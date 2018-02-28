# Kubernetes capstan

[![Build Status](https://travis-ci.org/ZJU-SEL/capstan.svg?branch=master)](https://travis-ci.org/ZJU-SEL/capstan)
[![codecov](https://codecov.io/gh/ZJU-SEL/capstan/branch/master/graph/badge.svg)](https://codecov.io/gh/ZJU-SEL/capstan)
[![Go Report Card](https://goreportcard.com/badge/github.com/ZJU-SEL/capstan)](https://goreportcard.com/report/github.com/ZJU-SEL/capstan)

## Introduction

capstan is a benchmarker which contains series of workloads and testing tools for Kubernetes. You can obtain the performance data of each type workload in the specific configuration of Kubernetes cluster offered by different cloud offerings.

## What is the scope of this project?

capstan aims to provide a series of workloads and testing tools for Kubernetes cluster:

- Run every workload's testing cases.

- Collect the testing results and the performance data of Kubernetes component and Kubernetes cluster.

- Analysis and display the testing results.

- Generate a testing report and performance report.

## What is not in scope for this project?

- Building a new cluster lifecycle management tool(e.g. [kubeadm](https://github.com/kubernetes/kubeadm),[kops](https://github.com/kubernetes/kops),[kubernetes-anywhere](https://github.com/kubernetes/kubernetes-anywhere)).

- Building a new data collection and analysis tool(e.g. [cadvisor](https://github.com/google/cadvisor),[heapster](https://github.com/kubernetes/heapster)).

## QuickStart

Build capstan:

```sh
mkdir -p $GOPATH/src/github.com/ZJU-SEL
git clone https://github.com/ZJU-SEL/capstan.git $GOPATH/src/github.com/ZJU-SEL/capstan
cd $GOPATH/src/github.com/ZJU-SEL/capstan
make && make install
```

Install Prometheus and Grafana:

```sh
# install Docker
apt-get install docker.io -y
# install Pushgateway
docker run -d -p 9091:9091 prom/pushgateway
# install Prometheus
cat >/tmp/prometheus.yml <<EOF
global:
  scrape_interval: 15s
  scrape_timeout: 10s
scrape_configs:
  - job_name: 'pushgateway'
    static_configs:
      - targets: ['<Your-HostIP>:9091']
EOF
docker run -d -p 9090:9090 -v /tmp/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus
# install Grafana
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
docker run -d -p 3000:3000 -v /tmp/provisioning/datasources:/etc/grafana/provisioning/datasources grafana/grafana
```

Prepare Kubernetes admin config file:

```sh
Copy the Kubernetes admin config file to the host path /etc/kubernetes/admin.conf
```

Configure capstan:

```sh
cat >/etc/capstan/config <<EOF
{
    "ResultsDir": "/tmp/capstan",
    "Prometheus": {
        "PushgatewayEndpoint": "http://<Your-HostIP>:9091"
    },
    "Steps": 10,
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

## Documentation

- [Deploying](docs/deploy.md)

## Roadmap

- Design the testing indicators（P0).
- Design and Implement the framework of capstan（P0).
- Use Prometheus and Grafana to analysis and display the testing results（P0).
- Implement multiple workloads and testing tools（P1).
- Implement the online ranking system（P2).
