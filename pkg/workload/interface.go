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

package workload

import (
	"github.com/ZJU-SEL/capstan/pkg/testingtool"
	"k8s.io/client-go/kubernetes"
)

// Interface should be implemented by a specific workload.
type Interface interface {
	// Run runs a testing workload.
	Run(kubeClient kubernetes.Interface) error
	// GetTestingTool returns the testing tool interface of this workload.
	GetTestingTool() testingtool.Interface
	// GetName returns the name of this workload.
	GetName() string
	// GetImage returns the image name of this workload.
	GetImage() string
}

// Workload is the internal representation of a testing workload.
type Workload struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	TestingTool testingtool.TestingTool
}

// DefWorkloads is the defined workloads.
var DefWorkloads = []string{
	"nginx",
	"iperf",
}
