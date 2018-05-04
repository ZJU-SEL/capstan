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

package helm

import (
	"time"

	"github.com/ZJU-SEL/capstan/pkg/workload"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

// Workload represents the workload managed by helm.
type Workload struct {
	workload  workload.Workload
	Name      string
	Helm      workload.Helm
	Frequency int
}

// Ensure the workload implements workload.Interface
var _ workload.Interface = &Workload{}

// NewWorkload creates a new workload from the given workload definition.
func NewWorkload(wl workload.Workload) (*Workload, error) {
	if wl.Helm.Name == "" && wl.Helm.Chart == "" {
		return nil, errors.Errorf("The helm.name and helm.chart of a workload must be not null")
	}
	return &Workload{
		workload:  wl,
		Name:      wl.Name,
		Helm:      wl.Helm,
		Frequency: wl.Frequency,
	}, nil
}

// Run runs the workload (to adhere to workload.Interface).
func (w *Workload) Run(kubeClient kubernetes.Interface) error {
	// initialize a new test tool for the workload.
	testTool, err := w.TestTool()
	if err != nil {
		return err
	}

	// make sure there is a clean namespace.
	err = workload.CleanNamespace(kubeClient, workload.Namespace)
	if err != nil {
		return errors.Wrapf(err, "Failed delete namespace %s", workload.Namespace)
	}

	for i := 1; i <= w.Frequency; i++ {
		for _, testCase := range testTool.GetTestCaseSet() {
			// running a test case.
			glog.V(1).Infof("Workload %s repeat %d: Running the test case %q of workload %s", w.Name, i, testCase.Name, w.Name)
			err := testTool.Run(kubeClient, testCase)
			if err != nil {
				_ = testTool.Cleanup(kubeClient)
				return errors.Wrapf(err, "Failed to create the resouces belong to test case %q of %s", testCase.Name, w.Name)
			}

			// check the test case has completed successfully.
			glog.V(4).Infof("Workload %s repeat %d: checking the test case %q has completed successfully", w.Name, i, testCase.Name)
			err = testTool.HasTestDone(kubeClient)
			if err != nil {
				_ = testTool.Cleanup(kubeClient)
				return errors.Wrapf(err, "The test case %s has not completed successfully", testCase.Name)
			}

			// clean up all the resouces created by the test case.
			glog.V(4).Infof("Workload %s repeat %d: Cleaning up all the resouces created by the test case %q", w.Name, i, testCase.Name)
			err = testTool.Cleanup(kubeClient)
			if err != nil {
				return errors.Wrapf(err, "Failed to cleanup the resouces created by the test case %s", testCase.Name)
			}

			// sleep some seconds between testing cases.
			glog.V(4).Infof("Workload %s repeat %d: Sleeping %v and starting next test case.", w.Name, i, testTool.GetSteps())
			time.Sleep(testTool.GetSteps())
		}
	}
	return nil
}

// TestTool initializes a new test tool for the workload (to adhere to workload.Interface).
func (w *Workload) TestTool() (workload.Tool, error) {
	return &TestTool{
		Workload:    w,
		Name:        w.workload.TestTool.Name,
		Script:      w.workload.TestTool.Script,
		Image:       w.workload.TestTool.Image,
		Steps:       time.Duration(w.workload.TestTool.Steps) * time.Second,
		TestCaseSet: w.workload.TestTool.TestCaseSet,
	}, nil
}
