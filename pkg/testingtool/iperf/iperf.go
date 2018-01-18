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

package iperf

import (
	"fmt"
	"strings"
	"time"

	"github.com/ZJU-SEL/capstan/pkg/testingtool"
	"k8s.io/client-go/kubernetes"
)

const (
	toolName = "iperf"
)

// TestingCaseSet is the list of wrk defined testing case.
var TestingCaseSet = []string{
	"benchmarkPodIPSameNode",
	"benchmarkVIPSameNode",
	"benchmarkPodIPDiffNode",
	"benchmarkVIPDiffNode",
}

// TestingTool represents the wrk testing tool.
type TestingTool struct {
	Name           string
	Image          string
	Steps          time.Duration
	TestingCaseSet []string
}

// Ensure iperf testing tool implements testingtool.Interface
var _ testingtool.Interface = &TestingTool{}

// NewTool creates a new iperf testing tool from the given testing tool definition.
func NewTool(tool testingtool.TestingTool) (*TestingTool, error) {
	if tool.Name != toolName {
		return nil, fmt.Errorf("Wrong parameter, the testing tool name must be %q", toolName)
	}

	return &TestingTool{
		Name:           toolName,
		Image:          tool.Image,
		Steps:          time.Duration(tool.Steps) * time.Second,
		TestingCaseSet: strings.Split(tool.TestingCaseSet, ","),
	}, nil
}

// Run runs the defined testing case set for iperf testing tool (to adhere to testingtool.Interface).
func (t *TestingTool) Run(kubeClient kubernetes.Interface, testingCase string) error {
	return fmt.Errorf("Not implemented")
}

//GetTestingResults gets the testing results of iperf testing case (to adhere to testingtool.Interface).
func (t *TestingTool) GetTestingResults(kubeClient kubernetes.Interface) error {
	return fmt.Errorf("Not implemented")
}

// Cleanup cleans up all resources created by a testing case for iperf testing tool (to adhere to testingtool.Interface).
func (t *TestingTool) Cleanup(kubeClient kubernetes.Interface) error {
	return fmt.Errorf("Not implemented")
}

// Monitor continually checks for problems in the resources created by a
// testing case (either because it won't schedule, too many failed executions, etc)
// and sends the errors through the provided channel (to adhere to testingtool.Interface).
func (t *TestingTool) Monitor(kubeClient kubernetes.Interface, testingErr chan error) {
	testingErr <- fmt.Errorf("Not implemented")
}

// GetName returns the name of iperf testing tool (to adhere to testingtool.Interface).
func (t *TestingTool) GetName() string {
	return t.Name
}

// GetImage returns the image name of iperf testing tool (to adhere to testingtool.Interface).
func (t *TestingTool) GetImage() string {
	return t.Image
}

// GetSteps returns the steps between each testing case (to adhere to testingtool.Interface).
func (t *TestingTool) GetSteps() time.Duration {
	return t.Steps
}

// GetTestingCaseSet returns the testing case set which the iperf testing tool will run (to adhere to testingtool.Interface).
func (t *TestingTool) GetTestingCaseSet() []string {
	return t.TestingCaseSet
}
