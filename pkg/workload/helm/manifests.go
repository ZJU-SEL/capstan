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

package helm

const (
	// PodAntiAffinity is the manifest for test tool running on the different node with workload.
	PodAntiAffinity = `
apiVersion: v1
kind: Pod
metadata:
  name: {{ .Name }}
  annotations:
    capstan-test: {{ .Name }}
    capstan-testcase: {{ .TestingName }}
  labels:
    component: capstan
    test: {{ .Name }}
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
  - name: {{ .Name }}
    image: {{ .Image }}
    imagePullPolicy: Always
    args: [{{ .Args }}]
    envFrom:
    - configMapRef:
        name: capstan-envs
    volumeMounts:
    - name: script-volume
      mountPath: /opt/capstan
  volumes:
    - name: script-volume
      configMap:
        name: capstan-script
        items:
        - key: run_test.sh
          path: run_test.sh
  restartPolicy: Never
  tolerations:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
    operator: Exists
  - key: CriticalAddonsOnly
    operator: Exists
  serviceAccountName: {{ .ServiceAccountName }}
`
	// PodAffinity is the manifest for test tool running on the same node with workload.
	PodAffinity = `
apiVersion: v1
kind: Pod
metadata:
  name: {{ .Name }}
  annotations:
    capstan-test: {{ .Name }}
    capstan-testcase: {{ .TestingName }}
  labels:
    component: capstan
    test: {{ .Name }}
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
  - name: {{ .Name }}
    image: {{ .Image }}
    imagePullPolicy: Always
    args: [{{ .Args }}]
    envFrom:
    - configMapRef:
        name: capstan-envs
    volumeMounts:
    - name: script-volume
      mountPath: /opt/capstan
  volumes:
    - name: script-volume
      configMap:
        name: capstan-script
        items:
        - key: run_test.sh
          path: run_test.sh
  restartPolicy: Never
  tolerations:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
    operator: Exists
  - key: CriticalAddonsOnly
    operator: Exists
  serviceAccountName: {{ .ServiceAccountName }}
`
	// PodAnyAffinity is the manifest for test tool running on the any node with workload.
	PodAnyAffinity = `
apiVersion: v1
kind: Pod
metadata:
  name: {{ .Name }}
  annotations:
    capstan-test: {{ .Name }}
    capstan-testcase: {{ .TestingName }}
  labels:
    component: capstan
    test: {{ .Name }}
  namespace: {{ .Namespace }}
spec:
  containers:
  - name: {{ .Name }}
    image: {{ .Image }}
    imagePullPolicy: Always
    args: [{{ .Args }}]
    envFrom:
    - configMapRef:
        name: capstan-envs
    volumeMounts:
    - name: script-volume
      mountPath: /opt/capstan
  volumes:
    - name: script-volume
      configMap:
        name: capstan-script
        items:
        - key: run_test.sh
          path: run_test.sh
  restartPolicy: Never
  tolerations:
  - effect: NoSchedule
    key: node-role.kubernetes.io/master
    operator: Exists
  - key: CriticalAddonsOnly
    operator: Exists
  serviceAccountName: {{ .ServiceAccountName }}
`
)
