#!/bin/sh
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

cat >cnn-benchmark.yaml <<-EOF
apiVersion: kubeflow.org/v1alpha1
kind: TFJob
metadata:
  name: cnn-benchmark
  namespace: capstan
spec:
  replicaSpecs:
  - replicas: 1
    template:
      spec:
        containers:
        - args:
          - python
          - tf_cnn_benchmarks.py
          - --batch_size=32
          - --model=resnet50
          - --variable_update=parameter_server
          - --flush_stdout=true
          - --num_gpus=1
          - --local_parameter_device=cpu
          - --device=cpu
          - --data_format=NHWC
          image: gcr.io/kubeflow/tf-benchmarks-cpu:v20171202-bdab599-dirty-284af3
          name: tensorflow
          workingDir: /opt/tf-benchmarks/scripts/tf_cnn_benchmarks
        restartPolicy: OnFailure
    tfReplicaType: WORKER
  - replicas: 1
    template:
      spec:
        containers:
        - args:
          - python
          - tf_cnn_benchmarks.py
          - --batch_size=32
          - --model=resnet50
          - --variable_update=parameter_server
          - --flush_stdout=true
          - --num_gpus=1
          - --local_parameter_device=cpu
          - --device=cpu
          - --data_format=NHWC
          image: gcr.io/kubeflow/tf-benchmarks-cpu:v20171202-bdab599-dirty-284af3
          name: tensorflow
          workingDir: /opt/tf-benchmarks/scripts/tf_cnn_benchmarks
        restartPolicy: OnFailure
    tfReplicaType: PS
  terminationPolicy:
    chief:
      replicaIndex: 0
      replicaName: WORKER
  tfImage: gcr.io/kubeflow/tf-benchmarks-cpu:v20171202-bdab599-dirty-284af3
EOF

kubectl create -f cnn-benchmark.yaml
if [ $? -ne 0 ]
then
    echo "Failed run kubeflow-bench $*"
    exit 1
fi

sleep 30

PODNAME=$(kubectl get pods -n capstan | grep cnn-benchmark-worker | awk '{print $1}')
if [ $? -ne 0 ]
then
    echo "Failed run kubeflow-bench $*"
    exit 1
fi

echo "The worker pod is "$PODNAME

sleep 60

while :
do
    FINISHED=$(kubectl logs $PODNAME -n capstan | grep Finished)
    if [ "$FINISHED" != "" ] 
    then
        break
    fi
    sleep 120
done

echo "Redirect the worker pod logs"
kubectl logs $PODNAME -n capstan > test-log 2>&1

if [ $? -ne 0 ]
then
    echo "Failed run kubeflow-bench $*"
    cat test-log
    exit 1
fi

# Resolve result
echo "Resolving result"
cat test-log

images_second=$(cat test-log | grep total | awk '{print $NF}')
echo $images_second

RESULT="images_second $images_second $PrometheusLabel"
echo $RESULT
