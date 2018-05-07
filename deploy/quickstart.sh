#!/bin/bash
# Copyright (c) 2018 The ZJU-SEL Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# config the dashboards
cd $GOPATH/src/github.com/ZJU-SEL/capstan/grafana-dashboards
./configDashboard.sh

# delete the conflicts
PORTS=(3000 9090 9091)
for((i=0;i<${#PORTS[@]};i++));do
    CONTAINERID=$(docker ps | grep ${PORTS[i]} | awk '{print $1}')
    if [ -n "$CONTAINERID" ]; then
        docker rm -f $CONTAINERID;
    fi
done

# install Pushgateway
docker run -d -p 9091:9091 prom/pushgateway

# install Node-exporter
if [ -n "$(kubectl get ns | grep  capstan-exporter)" ]; then
    kubectl delete ns  capstan-exporter
    printf 'deleting  capstan-exporter namespace'
    while [ -n "$(kubectl get ns | grep  capstan-exporter)" ]
    do
        printf '.'
    done
    printf '\n'
fi
sleep 5
kubectl create namespace capstan-exporter
kubectl create -f $GOPATH/src/github.com/ZJU-SEL/capstan/deploy/node-exporter.yaml

# Get user's HostIP
IPADDR=$(curl -s http://members.3322.org/dyndns/getip)
sed -i "s/<Your-HostIP>/${IPADDR}/g" /etc/capstan/grafana/provisioning/datasources/prometheus.yaml /etc/capstan/config

# Run Prometheus
docker run -d -p 9090:9090 -v /etc/capstan/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml prom/prometheus

# Run Grafana
docker run -d -p 3000:3000 -v /etc/capstan/grafana/provisioning:/etc/grafana/provisioning grafana/grafana
