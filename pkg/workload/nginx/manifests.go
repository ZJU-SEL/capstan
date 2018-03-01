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

package nginx

const (
	nginxPod = `
apiVersion: v1
kind: Pod
metadata:
  name: {{ .Name }}
  annotations:
    capstan-workload: nginx
    capstan-testingcase: {{ .TestingName }}
  labels:
    component: capstan
    testing: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  containers:
  - name: workload-nginx
    image: {{ .Image }}
    imagePullPolicy: Always
    ports:
    - containerPort: 80
  restartPolicy: Never
  tolerations:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
    operator: Exists
  - key: CriticalAddonsOnly
    operator: Exists
`
	wrkPodAntiAffinity = `
apiVersion: v1
kind: Pod
metadata:
  name: {{ .Name }}
  annotations:
    capstan-testing: wrk
    capstan-testingcase: {{ .TestingName }}
  labels:
    component: capstan
    testing: {{ .Name }}
  namespace: {{ .Namespace }}
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
  - name: testing-wrk
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
	wrkPodAffinity = `
apiVersion: v1
kind: Pod
metadata:
  name: {{ .Name }}
  annotations:
    capstan-testing: wrk
    capstan-testingcase: {{ .TestingName }}
  labels:
    component: capstan
    testing: {{ .Name }}
  namespace: {{ .Namespace }}
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
  - name: testing-wrk
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
