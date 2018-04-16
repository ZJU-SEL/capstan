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

	"k8s.io/client-go/kubernetes"
)

var (
	// Namespace is the namespace of Kubernetes where capstan creates resources.
	Namespace = "capstan"
)

// Interface should be implemented by a specific workload.
type Interface interface {
	// Run runs a testing workload.
	Run(kubeClient kubernetes.Interface) error
	// TestTool returns the workload‘s Tool interface.
	TestTool() (Tool, error)
}

// Tool should be implemented by a testing tool.
type Tool interface {
	// Run runs the defined test case set.
	Run(kubeClient kubernetes.Interface, testCase TestCase) error
	// HasTestDone check a test case has finish or not.
	HasTestDone(kubeClient kubernetes.Interface) error
	// Cleanup cleans up all resources created by a test case.
	Cleanup(kubeClient kubernetes.Interface) error
	// GetName returns the name of this test tool.
	GetName() string
	// GetSteps returns the steps between each test case.
	GetSteps() time.Duration
	// GetTestCaseSet returns the test case set which the test tool will run.
	GetTestCaseSet() []TestCase
	// GetWorkload returns a workload that the test tool to be test.
	GetWorkload() Workload
}

// Workload is the internal representation of a testing workload.
type Workload struct {
	Name      string `json:"name"`
	Helm      Helm
	Frequency int `json:"frequency"`
	TestTool  TestTool
}

// Helm is the internal representation of helm.
type Helm struct {
	Name  string `json:"name"`
	Set   string `json:"set"`
	Chart string `json:"chart"`
}

// TestTool is the internal representation of a test tool.
type TestTool struct {
	Name        string `json:"name"`
	Image       string `json:"image"`
	Script      string `json:"script"`
	Steps       int    `json:"steps"`
	TestCaseSet []TestCase
}

// TestCase is the internal representation of a test case.
type TestCase struct {
	Name     string `json:"name"`
	Affinity bool   `json:"affinity"`
	Args     string `json:"args"`
	Envs     string `json:"envs"`
	Metrics  string `json:"metrics"`
}
