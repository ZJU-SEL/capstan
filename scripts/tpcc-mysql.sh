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

# Install test tool
git clone https://github.com/Percona-Lab/tpcc-mysql.git \
    && cd tpcc-mysql/src \
    && make all \
    && cd ..

# Configration
MYSQL=/usr/bin/mysql
TPCCLOAD=./tpcc_load
TPCCSTART=./tpcc_start
TABLESQL=./create_table.sql
CONSTRAINTSQL=./add_fkey_idx.sql

DATABASE=tpcc
USER=root

check_result(){
    if [ $? -ne 0 ]
    then
        exit 1
    fi
} 

# Load data
echo 'Start Load data ...'
echo MYSQL_HOST=$MYSQL_HOST
echo MYSQL_PASSWORD=$MYSQL_ROOT_PASSWORD

echo 'Drop tpcc database if exist'
$MYSQL -h $MYSQL_HOST -P3306 -u $USER -p$MYSQL_ROOT_PASSWORD -e "DROP DATABASE IF EXISTS $DATABASE"
check_result

echo 'Create tpcc database'
$MYSQL -h $MYSQL_HOST -P3306 -u $USER -p$MYSQL_ROOT_PASSWORD -e "CREATE DATABASE $DATABASE"
check_result

echo 'create test table in tpcc database'
$MYSQL -h $MYSQL_HOST -P3306 -u $USER -p$MYSQL_ROOT_PASSWORD $DATABASE < $TABLESQL
check_result

echo 'create indexes and FK'
$MYSQL -h $MYSQL_HOST -P3306 -u $USER -p$MYSQL_ROOT_PASSWORD $DATABASE < $CONSTRAINTSQL
check_result

# Populate data
echo 'Populate data ...'
echo warehouses=$1
$TPCCLOAD -h$MYSQL_HOST -P3306 -d $DATABASE -u $USER -p$MYSQL_ROOT_PASSWORD $1
check_result

# Start benchmark
echo 'Start TPCC benchmark ...'
$TPCCSTART -h$MYSQL_HOST -P3306 -d $DATABASE -u $USER -p$MYSQL_ROOT_PASSWORD $* >test-log 2>&1
if [ $? -ne 0 ]
then
    echo "Failed run tpccstart $*"
    cat test-log
    exit 1
fi

# Resolve result
echo "Resolving result"
cat test-log

line=`grep " TpmC" test-log`

tpmc=`echo $line|cut -d " " -f1`

RESULT="TPMC $tpmc $PrometheusLabel"


