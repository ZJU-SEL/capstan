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
	"time"

	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

const (
	// DefaultNamespace is the default namespace for capstan.
	DefaultNamespace = "capstan"
)

// Interface should be implemented by a specific workload.
type Interface interface {
	// Run runs a testing workload.
	Run(kubeClient kubernetes.Interface) error
	// TestingTool returns the workload‘s Tool interface.
	TestingTool() (Tool, error)
	// GetName returns the name of this workload.
	GetName() string
	// GetImage returns the image name of this workload.
	GetImage() string
}

// Tool should be implemented by a testing tool.
type Tool interface {
	// Run runs the defined testing case set.
	Run(kubeClient kubernetes.Interface, testingCase TestingCase) error
	// GetTestingResults gets the testing results of a testing case.
	GetTestingResults(kubeClient kubernetes.Interface) error
	// Cleanup cleans up all resources created by a testing case.
	Cleanup(kubeClient kubernetes.Interface) error
	// GetName returns the name of this testing tool.
	GetName() string
	// GetImage returns the image name of this testing tool.
	GetImage() string
	// GetSteps returns the steps between each testing case.
	GetSteps() time.Duration
	// GetTestingCaseSet returns the testing case set which the testing tool will run.
	GetTestingCaseSet() []TestingCase
}

// Workload is the internal representation of a testing workload.
type Workload struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	Frequency   int    `json:"frequency"`
	TestingTool TestingTool
}

// TestingTool is the internal representation of a testing tool.
type TestingTool struct {
	Name           string `json:"name"`
	Image          string `json:"image"`
	Steps          int    `json:"Steps"`
	TestingCaseSet []TestingCase
}

// TestingCase is the internal representation of a testing case.
type TestingCase struct {
	Name            string `json:"name"`
	WorkloadArgs    string `json:"workloadArgs"`
	TestingToolArgs string `json:"testingToolArgs"`
}

// DefWorkloads is the defined workloads.
var DefWorkloads = []string{
	"nginx",
	"iperf",
}

// DefTools is list of the defined testing tools.
var DefTools = []string{
	"wrk",
	"iperf",
}

// TestingCaseSetHasDefined finds whether all the string in slice a have defined in slice b or not.
func TestingCaseSetHasDefined(testingCaseSet []TestingCase, defs []string) error {
	for _, testingCase := range testingCaseSet {
		found := false
		for _, def := range defs {
			if testingCase.Name == def {
				found = true
			}
		}
		if !found {
			return errors.Errorf("Testing case %v has not defined, the testingCaseSet must in %v", testingCase.Name, defs)
		}
	}
	return nil
}
