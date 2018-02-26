/*
Copyright (c) 2018 The ZJU-SEL Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package iperf3

const (
	iperfServerPod = `
apiVersion: v1
kind: Pod
metadata:
  name: {{ .Name }}
  annotations:
    capstan-workload: iperf3
    capstan-testingcase: {{ .TestingName }}
  labels:
    component: capstan
    testing: {{ .Name }}
  namespace: capstan
spec:
  containers:
  - name: workload-iperf3
    image: {{ .Image }}
    imagePullPolicy: Always
    args: ["-c"]
    ports:
    - containerPort: 5201
  restartPolicy: Never
  tolerations:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
    operator: Exists
  - key: CriticalAddonsOnly
    operator: Exists
`
	iperfClientPodAntiAffinity = `
apiVersion: v1
kind: Pod
metadata:
  name: {{ .Name }}
  annotations:
    capstan-testing: iperf3
    capstan-testingcase: {{ .TestingName }}
  labels:
    component: capstan
    testing: {{ .Name }}
  namespace: capstan
spec:
  affinity:
    podAntiAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
          - key: testing
            operator: In
            values:
            -  {{ .WorkloadName }}
        topologyKey: "kubernetes.io/hostname"
  containers:
  - name: testing-iperf3
    image: {{ .Image }}
    imagePullPolicy: Always
    args: [{{ .Args }}]
    env:
    - name: ENDPOINT
      value: {{ .PodIP }}
  restartPolicy: Never
  tolerations:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
    operator: Exists
  - key: CriticalAddonsOnly
    operator: Exists
`
	iperfClientPodAffinity = `
apiVersion: v1
kind: Pod
metadata:
  name: {{ .Name }}
  annotations:
    capstan-testing: iperf3
    capstan-testingcase: {{ .TestingName }}
  labels:
    component: capstan
    testing: {{ .Name }}
  namespace: capstan
spec:
  affinity:
    podAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
      - labelSelector:
          matchExpressions:
          - key: testing
            operator: In
            values:
            -  {{ .WorkloadName }}
        topologyKey: "kubernetes.io/hostname"
  containers:
  - name: testing-iperf3
    image: {{ .Image }}
    imagePullPolicy: Always
    args: [{{ .Args }}]
    env:
    - name: ENDPOINT
      value: {{ .PodIP }}
  restartPolicy: Never
  tolerations:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
    operator: Exists
  - key: CriticalAddonsOnly
    operator: Exists
`
)
