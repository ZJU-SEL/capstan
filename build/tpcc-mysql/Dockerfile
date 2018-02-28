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

FROM debian:jessie

MAINTAINER The ZJU-SEL team

RUN apt-get update && apt-get install -y --force-yes --no-install-recommends \
    build-essential \
    ca-certificates \
    libmysql++-dev \
    git

RUN apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

RUN git clone https://github.com/Percona-Lab/tpcc-mysql.git /opt/tpcc-mysql \
    && cd /opt/tpcc-mysql/src \
    && make all

ENV PATH /opt/tpcc-mysql:$PATH
WORKDIR /opt/tpcc-mysql

COPY run_tpcc_mysql.sh /run_tpcc_mysql.sh 
ENTRYPOINT ["/run_tpcc_mysql.sh"]
