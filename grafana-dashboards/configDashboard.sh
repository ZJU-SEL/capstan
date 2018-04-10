#!/bin/bash
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

# set variables
CONFIGFILE="/etc/capstan/config"
SINGLESTATFILE="$GOPATH/src/github.com/ZJU-SEL/capstan/grafana-dashboards/singlestat.json"
EXPRFILE="$GOPATH/src/github.com/ZJU-SEL/capstan/grafana-dashboards/expr.json"
PATTERNFILE="$GOPATH/src/github.com/ZJU-SEL/capstan/grafana-dashboards/pattern.json"
SINGLESTAT_ALL_LINES=81
EXPR_ALL_LINES=8
PATTERN_ALL_LINES=18
SINGLESTAT_BASIC_LINE=300
EXPR_BASIC_LINE=214
PATTERN_BASIC_LINE=204

REFID=(A B C D E F)

declare -A NGINX=(["QPS"]="capstan_QPS")
declare -A IPERF3=(["BandWidth"]="capstan_BandWidth")
declare -A MYSQL=(["TPMC"]="capstan_TPMC")

PROMQL=""

# Get the PromQL according to workloadName and metrics
function GetPromQL(){
	case $1 in
	"nginx")
		PROMQL=${NGINX[$2]}
		;;
	"iperf3")
		PROMQL=${IPERF3[$2]}
		;;
	"mysql")
		PROMQL=${MYSQL[$2]}
		;;
	*)
		# TODO(ZeroMagic): address the non-existent workloads
		echo error
		;;
	esac
	return 0
}

# start config
WORKLOADSNUM=$(cat ${CONFIGFILE} | jq ".Workloads | length")

for((i=0;i<${WORKLOADSNUM};i++));do
	
	# get WorkloadName
	WORKLOADNAME=$(cat ${CONFIGFILE} | jq ".Workloads[${i}].name" | sed -e 's/"//g')

	cp $GOPATH/src/github.com/ZJU-SEL/capstan/grafana-dashboards/root.json /etc/capstan/grafana/provisioning/dashboards/${WORKLOADNAME}.json
	DESTFILE="/etc/capstan/grafana/provisioning/dashboards/${WORKLOADNAME}.json"

	# record all unique metrics which appear in the workload
	METRICSALL=()
	TOT=0
	
	# record the number of singlestats 
	COUNT=0

	# get testCase
	TESTCASESNUM=$(cat ${CONFIGFILE} | jq ".Workloads[${i}].testTool.testCaseSet | length")
	
	for((j=0;j<${TESTCASESNUM};j++));do
		TESTCASENAME=$(cat ${CONFIGFILE} | jq ".Workloads[${i}].testTool.testCaseSet[${j}].name" | sed -e 's/"//g')

		# get metrics
		ALLMETRICS=$(cat ${CONFIGFILE} | jq ".Workloads[${i}].testTool.testCaseSet[${j}].metrics" | sed -e 's/"//g')
		METRICS=(${ALLMETRICS//,/ })
		METRICSNUM=${#METRICS[@]}
		for((k=0;k<${METRICSNUM};k++));do
			((COUNT++))
			
			# judge the current metrics whether has appeared in this workload
			flag="true"
			for((x=0;x<${#METRICSALL[@]};x++));do
				if [ ${METRICSALL[${x}]} == ${METRICS[${k}]} ]
				then
					flag="false"
					break
				fi
			done
			
			if [ $flag == "true" ]
			then
				METRICSALL[${TOT}]=${METRICS[${k}]}
				((TOT++))
				if [ ${TOT} -eq 1 ]
				then
					# modify metrics
					sed -i "s/#metrics/${METRICS[${k}]}/g" ${DESTFILE} 
				else
					if [ ${#METRICSALL[@]} -eq 2 ]
					then
						sed -i "197s/Value/Value #A/g" ${DESTFILE}
					fi

					# calculate the line where the pattern should be added.
					PATTERN_NEWLINE=$[PATTERN_BASIC_LINE+PATTERN_ALL_LINES*(TOT-2)]
					sed -i "${PATTERN_NEWLINE} r ${PATTERNFILE}" ${DESTFILE}
					sed -i "$[PATTERN_NEWLINE+2]s/#metrics/${METRICS[${k}]}/g" ${DESTFILE}
					sed -i "$[PATTERN_NEWLINE+11]s/#value/Value #${REFID[$[TOT-1]]}/g" ${DESTFILE}

					# calculate the line where the expr should be added.
					EXPR_NEWLINE=$[EXPR_BASIC_LINE+EXPR_ALL_LINES*(TOT-2)+PATTERN_ALL_LINES*(TOT-1)]
					sed -i "${EXPR_NEWLINE} r ${EXPRFILE}" ${DESTFILE}
					sed -i "$[EXPR_NEWLINE+7]s/#A/${REFID[$[TOT-1]]}/g" ${DESTFILE}
				fi
			fi
			
			if [ ${COUNT} -gt 1 ]
			then
				# calculate the line where the singlestat should be added.
				SINGLESTAT_NEWLINE=$[SINGLESTAT_BASIC_LINE+SINGLESTAT_ALL_LINES*(COUNT-2)+EXPR_ALL_LINES*(TOT-1)+PATTERN_ALL_LINES*(TOT-1)]
				sed -i "${SINGLESTAT_NEWLINE} r ${SINGLESTATFILE}" ${DESTFILE}
				
				# set id
				sed -i "$[SINGLESTAT_NEWLINE+25]s/null/$[13+COUNT]/g" ${DESTFILE}
			fi
			
			# modify workloadName
			sed -i "s/#workloadName/${WORKLOADNAME}/g" ${DESTFILE}
			
			# modify PromQL
			GetPromQL ${WORKLOADNAME} ${METRICS[${k}]}
			sed -i "s/#PromQL/${PROMQL}/g" ${DESTFILE}
			
			# modify testCase
			LINE=$(grep -n "#testCase" ${DESTFILE} | awk '{print $1}' | sed -e 's/://g')
			sed -i "${LINE}s/#testCase/${TESTCASENAME}/g" ${DESTFILE}
			LINE=$(grep -n "#MetricsTestCase" ${DESTFILE} | awk '{print $1}' | sed -e 's/://g')
			sed -i "${LINE}s/#MetricsTestCase/${TESTCASENAME}-${METRICS[${k}]}/g" ${DESTFILE}
		done
	done

	WIDTH=$[24/${COUNT}]
	for((j=0;j<${COUNT};j++));do
		LINE=$(grep -n "\"w\"" ${DESTFILE} | sed -n "$[3+j], 1p" | awk '{print $1}' | sed -e 's/://g')

		# set width
		sed -i "${LINE}s/4/${WIDTH}/g" ${DESTFILE}
		# set x-coordinate
		sed -i "$[LINE+1]s/0/$[WIDTH*${j}]/g" ${DESTFILE}
	done
	
done

