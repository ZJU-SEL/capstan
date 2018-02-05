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
	"time"

	"github.com/ZJU-SEL/capstan/pkg/workload"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

// Workload represents the nginx workload.
type Workload struct {
	workload  workload.Workload
	Name      string
	Image     string
	Frequency int
}

// Ensure nginx Workload implements workload.Interface
var _ workload.Interface = &Workload{}

// NewWorkload creates a new nginx workload from the given workload definition.
func NewWorkload(wl workload.Workload) *Workload {
	return &Workload{
		workload:  wl,
		Name:      wl.Name,
		Image:     wl.Image,
		Frequency: wl.Frequency,
	}
}

// Run runs a nginx workload (to adhere to workload.Interface).
func (w *Workload) Run(kubeClient kubernetes.Interface) error {
	// initialize a new testing tool for this nginx workload.
	testingTool, err := w.TestingTool()
	if err != nil {
		return err
	}

	for i := 1; i <= w.Frequency; i++ {
		for _, testingCase := range testingTool.GetTestingCaseSet() {
			// running a testing case.
			glog.V(1).Infof("Repeat %d: Running the testing case %q of %s", i, testingCase.Name, w.GetName())
			err := testingTool.Run(kubeClient, testingCase)
			if err != nil {
				return errors.Wrapf(err, "Failed to create the resouces belong to testing case %q of %s", testingCase.Name, w.GetName())
			}

			// get the testing results of the testing case.
			glog.V(4).Infof("Repeat %d: Starting fetch the testing results of the testing case %q", i, testingCase.Name)
			err = testingTool.GetTestingResults(kubeClient)
			if err != nil {
				return errors.Wrapf(err, "Failed to gets the testing results of the testing case %s", testingCase.Name)
			}

			// clean up all the resouces created by the testing case.
			glog.V(4).Infof("Repeat %d: Cleaning up all the resouces created by the testing case %q", i, testingCase.Name)
			err = testingTool.Cleanup(kubeClient)
			if err != nil {
				return errors.Wrapf(err, "Failed to cleanup the resouces created by the testing case %s", testingCase.Name)
			}

			// sleep some seconds between testing cases.
			glog.V(4).Infof("Repeat %d: Sleeping %v and starting next testing case.", i, testingTool.GetSteps())
			time.Sleep(testingTool.GetSteps())
		}
	}
	return nil
}

// TestingTool initializes a new testing tool for this nginx workload (to adhere to workload.Interface).
func (w *Workload) TestingTool() (workload.Tool, error) {
	// TODO(mozhuli): support one workload mapping many testing tools.
	if w.workload.TestingTool.Name != toolName {
		return nil, errors.Errorf("Wrong parameter(%q), the testing tool name must be %q", w.workload.TestingTool.Name, toolName)
	}

	if err := workload.TestingCaseSetHasDefined(w.workload.TestingTool.TestingCaseSet, TestingCaseSet); err != nil {
		return nil, err
	}

	return &TestingTool{
		Workload:       w,
		Name:           toolName,
		Image:          w.workload.TestingTool.Image,
		Steps:          time.Duration(w.workload.TestingTool.Steps) * time.Second,
		TestingCaseSet: w.workload.TestingTool.TestingCaseSet,
	}, nil
}

// GetName returns the name of this nginx workload (to adhere to workload.Interface).
func (w *Workload) GetName() string {
	return w.Name
}

// GetImage returns the image name of this nginx workload (to adhere to workload.Interface).
func (w *Workload) GetImage() string {
	return w.Image
}
