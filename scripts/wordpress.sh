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

#get the default blog address
year=`date +%Y`
month=`date +%m`
day=`date +%d`
other="/"
title="hello-world"
url=$*$other$year$other$month$other$day$other$title$other

echo $url

/usr/local/bin/wrk $url > test-log 2>&1
if [ $? -ne 0 ]
then
    echo "Failed run wrk $*"
    cat test-log
    exit 1
fi

# Resolve result
echo "Resolving result"
cat test-log

line1=`grep "Requests/sec:" test-log`
qps=`echo $line1|cut -d " " -f2`
RESULT="QPS $qps $PrometheusLabel"
