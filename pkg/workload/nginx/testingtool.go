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

package nginx

import (
	"fmt"
	"time"

	"github.com/ZJU-SEL/capstan/pkg/workload"
	"k8s.io/client-go/kubernetes"
)

const (
	toolName = "wrk"
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
	Workload       *Workload
	Name           string
	Image          string
	Steps          time.Duration
	TestingCaseSet []workload.TestingCase
}

// Ensure wrk testing tool implements workload.Tool interface.
var _ workload.Tool = &TestingTool{}

// Run runs the defined testing case set for wrk testing tool (to adhere to workload.Tool interface).
func (t *TestingTool) Run(kubeClient kubernetes.Interface, testingCase string) error {
	return fmt.Errorf("Not implemented")
}

//GetTestingResults gets the testing results of wrk testing case (to adhere to workload.Tool interface).
func (t *TestingTool) GetTestingResults(kubeClient kubernetes.Interface) error {
	return fmt.Errorf("Not implemented")
}

// Cleanup cleans up all resources created by a testing case for wrk testing tool (to adhere to workload.Tool interface).
func (t *TestingTool) Cleanup(kubeClient kubernetes.Interface) error {
	return fmt.Errorf("Not implemented")
}

// Monitor continually checks for problems in the resources created by a
// testing case (either because it won't schedule, too many failed executions, etc)
// and sends the errors through the provided channel (to adhere to workload.Tool interface).
func (t *TestingTool) Monitor(kubeClient kubernetes.Interface, testingErr chan error) {
	testingErr <- fmt.Errorf("Not implemented")
}

// GetName returns the name of wrk testing tool (to adhere to workload.Tool interface).
func (t *TestingTool) GetName() string {
	return t.Name
}

// GetImage returns the image name of wrk testing tool (to adhere to workload.Tool interface).
func (t *TestingTool) GetImage() string {
	return t.Image
}

// GetSteps returns the steps between each testing case (to adhere to workload.Tool interface).
func (t *TestingTool) GetSteps() time.Duration {
	return t.Steps
}

// GetTestingCaseSet returns the testing case set which the wrk testing tool will run (to adhere to workload.Tool interface).
func (t *TestingTool) GetTestingCaseSet() []workload.TestingCase {
	return t.TestingCaseSet
}
