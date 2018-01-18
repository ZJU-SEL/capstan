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

package testingtool

import (
	"time"

	"k8s.io/client-go/kubernetes"
)

// Interface should be implemented by a testing tool.
type Interface interface {
	// Run runs the defined testing case set.
	Run(kubeClient kubernetes.Interface, testingCase string) error
	//GetTestingResults gets the testing results of a testing case.
	GetTestingResults(kubeClient kubernetes.Interface) error
	// Cleanup cleans up all resources created by a testing case.
	Cleanup(kubeClient kubernetes.Interface) error
	// Monitor continually checks for problems in the resources created by a
	// testing case (either because it won't schedule, too many failed executions, etc)
	// and sends the errors through the provided channel.
	Monitor(kubeClient kubernetes.Interface, testingErr chan error)
	// GetName returns the name of this testing tool.
	GetName() string
	// GetImage returns the image name of this testing tool.
	GetImage() string
	// GetSteps returns the steps between each testing case.
	GetSteps() time.Duration
	// GetTestingCaseSet returns the testing case set which the testing tool will run.
	GetTestingCaseSet() []string
}

// TestingTool is the internal representation of a testing tool.
type TestingTool struct {
	Name           string `json:"name"`
	Image          string `json:"image"`
	Steps          int    `json:"Steps"`
	TestingCaseSet string `json:"testingCaseSet"`
}

// DefTools is list of the defined testing tools.
var DefTools = []string{
	"wrk",
	"iperf",
}
