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

	"github.com/ZJU-SEL/capstan/pkg/testingtool"
	"github.com/ZJU-SEL/capstan/pkg/testingtool/wrk"
	"github.com/ZJU-SEL/capstan/pkg/workload"
	"github.com/golang/glog"
	"k8s.io/client-go/kubernetes"
)

// Workload represents the nginx workload.
type Workload struct {
	Name        string
	Image       string
	TestingTool testingtool.Interface
}

// Ensure nginx Workload implements workload.Interface
var _ workload.Interface = &Workload{}

// NewWorkload creates a new nginx workload from the given workload definition.
func NewWorkload(wl workload.Workload) (*Workload, error) {
	// TODO(mozhuli): support one workload mapping many testing tools.
	tt, err := wrk.NewTool(wl.TestingTool)
	if err != nil {
		return nil, err
	}

	return &Workload{
		Name:        wl.Name,
		Image:       wl.Image,
		TestingTool: tt,
	}, nil
}

// Run runs a nginx workload (to adhere to workload.Interface).
func (w *Workload) Run(kubeClient kubernetes.Interface) error {
	for i, testingCase := range w.TestingTool.GetTestingCaseSet() {
		// running a testing case.
		glog.V(1).Infof("Running the %dth testing case:%v", i, testingCase)
		err := w.TestingTool.Run(kubeClient, testingCase)
		if err != nil {
			return fmt.Errorf("Failed to create the resouces belong to testing case %q :%v", testingCase, err)
		}

		// monitor the process of the testing case.
		testingErr := make(chan error)

		go w.TestingTool.Monitor(kubeClient, testingErr)

		err = <-testingErr
		if err != nil {
			return fmt.Errorf("Failed to test the case %q :%v", testingCase, err)
		}

		// get the testing results of the testing case.
		err = w.TestingTool.GetTestingResults(kubeClient)
		if err != nil {
			return fmt.Errorf("Failed to gets the testing results of the testing case %q :%v", testingCase, err)
		}

		// clean up the resouces created by the testing case.
		err = w.TestingTool.Cleanup(kubeClient)
		if err != nil {
			return fmt.Errorf("Failed to cleanup the resouces created by the testing case %q :%v", testingCase, err)
		}

		// sleep some seconds between testing cases.
		time.Sleep(w.TestingTool.GetSteps())
	}
	return nil
}

// GetTestingTool returns the testing tool interface for this nginx workload (to adhere to workload.Interface).
func (w *Workload) GetTestingTool() testingtool.Interface {
	return w.TestingTool
}

// GetName returns the name of this nginx workload (to adhere to workload.Interface).
func (w *Workload) GetName() string {
	return w.Name
}

// GetImage returns the image name of this nginx workload (to adhere to workload.Interface).
func (w *Workload) GetImage() string {
	return w.Name
}
