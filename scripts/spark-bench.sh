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


wget https://github.com/CODAIT/spark-bench/releases/download/v99/spark-bench_2.3.0_0.4.0-RELEASE_99.tgz \
    && tar -xf spark-bench_2.3.0_0.4.0-RELEASE_99.tgz \
    && rm spark-bench_2.3.0_0.4.0-RELEASE_99.tgz \
    && cd spark-bench_2.3.0_0.4.0-RELEASE

./bin/spark-bench.sh ./examples/minimal-example.conf > test-log 2>&1

if [ $? -ne 0 ]
then
    echo "Failed run spark-bench $*"
    cat test-log
    exit 1
fi

# Resolve result
echo "Resolving result"
cat test-log

line=`grep "|sparkpi|" test-log`

total_runtime=`echo $line|cut -d "|" -f4|cut -d " " -f2`

RESULT="total_runtime $total_runtime $PrometheusLabel"
