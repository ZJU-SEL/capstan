# How to add a new workload and test it

- [Overview](#overview)
- [Prerequisites](#prerequisites)
- [Add a new workload](#add-a-new-workload)
  - [Offer a chart](#offer-a-chart)
  - [Write a test script](#write-a-test-script)
  - [Configure capstan](#configure-capstan)

## Overview

This document shows how to add a new workload and use a test script to test it.

## Prerequisites

You have installed the capstan, if not, please follow the [deploying document](deploy.md)

## Add a new workload

### Offer a chart

We only support workload which managed by helm for now. So you should have a standard chart:

1. The workload should has a lable `app: {{ template "<Your-chart-name>.fullname" . }}`.

1. The workload should has a service and the service name must be `name: {{ template "<Your-chart-name>.fullname" . }}`.

### Write a test script

Secondly, we should have a test script which has the following required:

1. You should check each command's exit code. The test script should exit with non-zero in time once there is a non-zero exit code.

1. The test result should has the format:

   ```sh
   # for no more label
   RESULT="<metrics-name> <metrics-data> $PrometheusLabel"

   # add your own labels
   YourLabels="key1=value1,key2=value2"
   RESULT="<metrics-name> <metrics-data> $PrometheusLabel+$YourLabels"

   # more than one metrics
   YourLabels1="key1=value1,key2=value2"
   YourLabels2="key1=value1,key2=value2"
   RESULT="<metrics-name1> <metrics-data1> $PrometheusLabel+$YourLabels1\n<metrics-name2> <metrics-data2> $PrometheusLabel+$YourLabels2"
   ```

### Configure capstan

You should add some configuration in the `workload section` of capstan configuration file.

```json
        {
            "name": "<Your-chart-name>",
            "helm": {
                "name": "<Your-workload-release-name>",
                "set": "<Your-chart-set>",
                "chart": "<Your-chart-package-dir>"
            },
            "frequency": 1,
            "testTool": {
                "name": "<Your-test-name>",
                "script": "<Your-test-script-dir>",
                "steps": 10,
                "testCaseSet": [
                    {
                        "name": "<Your-test-case-name>",
                        "affinity": true,
                        "args": "<Your-test-script-agrs>",
                        "envs": "<Your-test-script-envs>",
                        "metrics": "<Your-metrics>"
                    },
                    {
                        "name": "<Your-test-case-name>",
                        "affinity": false,
                        "args": "<Your-test-script-agrs>",
                        "envs": "<Your-test-script-envs>",
                        "metrics": "<Your-metrics>"
                    }
                    // more test case ...
                ]
            }
        }
```

For more examples, please see the charts, scripts and examples/capstan.conf in the repository.
