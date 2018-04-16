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

# Ensure all variables are defined.
set -u

echo "The benchmark variables: $*"

# The run_benchmark.sh script is the default entrypoint, 
# calling run_test.sh which specified by user.
echo "Calling run_test.sh"
cp /opt/capstan/run_test.sh /root/capstan/run_test.sh
source ./run_test.sh $*

# if run_test.sh exist code not equal 0, exit 1.
if [ $? -ne 0 ]
then
    echo "Failed run ./run_test.sh $*"
    exit 1
fi

# Calling capstan-pusher to push the result to pushGateway.
echo "Calling capstan-pusher to push the result to pushGateway"

/root/capstan/capstan-pusher --endpoint=$PushgatewayEndpoint "$RESULT"
if [ $? -ne 0 ]
then
    echo "Failed to push result to pushGateway"
    exit 1
fi

echo "Capstan finish the test case"

