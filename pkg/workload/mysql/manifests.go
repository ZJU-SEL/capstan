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

package mysql

const (
	mysqlTPCCPodAntiAffinity = `
apiVersion: v1
kind: Pod
metadata:
  name: {{ .Name }}
  annotations:
    capstan-testing: mysql-tpcc
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
          - key: app
            operator: In
            values:
            -  {{ .Label }}
        topologyKey: "kubernetes.io/hostname"
  containers:
  - name: testing-mysql-tpcc
    image: {{ .Image }}
    imagePullPolicy: Always
    args: [{{ .Args }}]
    env:
    - name: MYSQL_HOST
      value: {{ .DNSName }}
    - name: MYSQL_ROOT_PASSWORD
      value: {{ .PASSWORD }}
  restartPolicy: Never
  tolerations:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
    operator: Exists
  - key: CriticalAddonsOnly
    operator: Exists
`
	mysqlTPCCPodAffinity = `
apiVersion: v1
kind: Pod
metadata:
  name: {{ .Name }}
  annotations:
    capstan-testing: mysql-tpcc
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
          - key: app
            operator: In
            values:
            -  {{ .Label }}
        topologyKey: "kubernetes.io/hostname"
  containers:
  - name: testing-mysql-tpcc
    image: {{ .Image }}
    imagePullPolicy: Always
    args: [{{ .Args }}]
    env:
    - name: MYSQL_HOST
      value: {{ .DNSName }}
    - name: MYSQL_ROOT_PASSWORD
      value: {{ .PASSWORD }}
  restartPolicy: Never
  tolerations:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
    operator: Exists
  - key: CriticalAddonsOnly
    operator: Exists
`
)
