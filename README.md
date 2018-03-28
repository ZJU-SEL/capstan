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

## Prerequisites Installation

- Go: 1.9.3+
- Helm

## QuickStart

In the quickstart, we use the default config to run capstan. You can also specify your own config.

### Prepare Kubernetes admin config file:

```sh
Copy the Kubernetes admin config file to the host path /etc/kubernetes/admin.conf

export KUBECONFIG=/etc/kubernetes/admin.conf
```

### Install kubectl:
```sh
curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl

chmod +x ./kubectl

sudo mv ./kubectl /usr/local/bin/kubectl
```

### Install helm

Here is [helm installtion](https://github.com/kubernetes/helm/blob/master/docs/install.md?spm=a2c4g.11186623.2.7.qwiWKY&file=install.md).
After installtion, you should `helm init`.

### Build capstan:

```sh
mkdir -p $GOPATH/src/github.com/ZJU-SEL
git clone https://github.com/ZJU-SEL/capstan.git $GOPATH/src/github.com/ZJU-SEL/capstan
cd $GOPATH/src/github.com/ZJU-SEL/capstan
make && make install
```

### Deploy capstan:

```sh
# install Docker
apt-get install docker.io -y

# configure Prometheus
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
  - job_name: 'gcp'
    static_configs:
      - targets: ['<Your-Kubernetes in gcp-MasterIP>:31672','<Your-Kubernetes in gcp-Node1IP>:31672','<Your-Kubernetes in gcp-Node2IP>:31672',...]
EOF

# deploy
sh $GOPATH/src/github.com/ZJU-SEL/capstan/deploy/quickstart.sh

```

### Start capstan:

```sh
capstan --v=3 --logtostderr --config=/etc/capstan/config --kubeconfig=/etc/kubernetes/admin.conf &
```

### Display

You can visit `<Your-HostIP>:3000` to see Grafana. There is a default user "admin", and its password is "admin". 

## Documentation

- [Deploying](docs/deploy.md)

## Roadmap

- Design the testing indicators（P0).
- Design and Implement the framework of capstan（P0).
- Use Prometheus and Grafana to analysis and display the testing results（P0).
- Implement multiple workloads and testing tools（P1).
- Implement the online ranking system（P2).