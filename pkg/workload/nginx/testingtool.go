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
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/ZJU-SEL/capstan/pkg/capstan/types"
	"github.com/ZJU-SEL/capstan/pkg/workload"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	toolName               = "wrk"
	benchmarkPodIPSameNode = "benchmarkPodIPSameNode"
	benchmarkPodIPDiffNode = "benchmarkPodIPDiffNode"
)

// TestingCaseSet is the list of wrk defined testing case.
var TestingCaseSet = []string{
	"benchmarkPodIPSameNode",
	"benchmarkPodIPDiffNode",
}

// TestingTool represents the wrk testing tool.
type TestingTool struct {
	Workload       *Workload
	Name           string
	Image          string
	Steps          time.Duration
	CurrentTesting workload.TestingCase
	TestingCaseSet []workload.TestingCase
}

// Ensure wrk testing tool implements workload.Tool interface.
var _ workload.Tool = &TestingTool{}

// Run runs the defined testing case set for wrk testing tool (to adhere to workload.Tool interface).
func (t *TestingTool) Run(kubeClient kubernetes.Interface, testingCase workload.TestingCase) error {
	t.CurrentTesting = testingCase

	// 1. start a workload for the testing case.
	workloadPodName := t.workloadPodName(testingCase.Name)
	tempWorkloadArgs := struct{ Name, TestingName, Image string }{
		Name:        workloadPodName,
		TestingName: testingCase.Name,
		Image:       t.Workload.GetImage(),
	}

	nginxPodBytes, err := workload.ParseTemplate(nginxPod, tempWorkloadArgs)
	if err != nil {
		return errors.Wrapf(err, "unable to parse %v using %v", nginxPod, tempWorkloadArgs)
	}

	glog.V(4).Infof("Creating workload %q of testing case %s", workloadPodName, testingCase.Name)
	if err := workload.CreatePod(kubeClient, nginxPodBytes); err != nil {
		return errors.Wrapf(err, "unable to create the %s workload for testing case %s", t.Workload.GetName(), testingCase.Name)
	}

	// 2. get the pod ip of the workload until workload is running.
	glog.V(4).Infof("Geting the pod ip of workload %s", workloadPodName)
	podIP, err := t.getRunningPodIP(kubeClient, workloadPodName)
	if err != nil {
		return errors.Wrapf(err, "unable to find pod %s's ip created by the %s workload for testing case %s", workloadPodName, t.Workload.GetName(), testingCase.Name)
	}

	// 3. start a testing pod for testing the workload.
	testingPodName := t.testingPodName(testingCase.Name)
	testingPod, args := t.findTemplate(testingCase.Name)
	tempTestingArgs := struct{ Name, TestingName, Image, WorkloadName, Args, PodIP string }{
		Name:         testingPodName,
		TestingName:  testingCase.Name,
		Image:        t.GetImage(),
		WorkloadName: workloadPodName,
		Args:         workload.FomatArgs(args),
		PodIP:        podIP,
	}

	testingPodBytes, err := workload.ParseTemplate(testingPod, tempTestingArgs)
	if err != nil {
		return errors.Wrapf(err, "unable to parse %v using %v", testingPod, tempTestingArgs)
	}

	glog.V(4).Infof("Creating testing pod %q of testing case %s", testingPodName, testingCase.Name)
	if err := workload.CreatePod(kubeClient, testingPodBytes); err != nil {
		return errors.Wrapf(err, "unable to create the testing pod for testing case %s", testingCase.Name)
	}

	return nil
}

//GetTestingResults gets the testing results of wrk testing case (to adhere to workload.Tool interface).
func (t *TestingTool) GetTestingResults(kubeClient kubernetes.Interface) error {
	name := t.testingPodName(t.CurrentTesting.Name)
	for {
		// Sleep between each poll
		// TODO(mozhuli): Use a watcher instead of polling.
		time.Sleep(30 * time.Second)

		// Make sure there's a pod.
		pod, err := kubeClient.CoreV1().Pods(workload.DefaultNamespace).Get(name, apismetav1.GetOptions{})
		if err != nil {
			return errors.WithStack(err)
		}

		// Make sure the pod isn't failing.
		if isFailing, err := workload.IsPodFailing(pod); isFailing {
			return err
		}

		// Check testing has done.
		body, err := kubeClient.CoreV1().Pods(workload.DefaultNamespace).GetLogs(
			name,
			&v1.PodLogOptions{},
		).Do().Raw()

		if err != nil {
			return errors.WithStack(err)
		}

		glog.V(5).Infof("Checking testing has done:\n%s", string(body))
		if workload.HasTestingDone(body) {
			glog.V(4).Infof("Testing case %s has done", t.CurrentTesting.Name)
			outdir := path.Join(types.ResultsDir, "workloads", t.Workload.GetName(), t.GetName(), t.CurrentTesting.Name)
			if err = os.MkdirAll(outdir, 0755); err != nil {
				return errors.WithStack(err)
			}

			outfile := path.Join(outdir, t.GetName()) + ".log"
			if err = ioutil.WriteFile(outfile, body, 0644); err != nil {
				return errors.WithStack(err)
			}
			return nil
		}
	}
}

// Cleanup cleans up all resources created by a testing case for wrk testing tool (to adhere to workload.Tool interface).
func (t *TestingTool) Cleanup(kubeClient kubernetes.Interface) error {
	if err := workload.DeletePod(kubeClient, t.testingPodName(t.CurrentTesting.Name)); err != nil {
		return err
	}
	if err := workload.DeletePod(kubeClient, t.workloadPodName(t.CurrentTesting.Name)); err != nil {
		return err
	}
	return nil
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

// getRunningPodIP finds the running pod's ip created by a workload, If no pod is found,
// or if pod's status is not running, returns an error.
func (t *TestingTool) getRunningPodIP(kubeClient kubernetes.Interface, name string) (string, error) {
	n := 0
	for {
		// Sleep between each poll, which should give the workload enough time to create a Pod
		// TODO(mozhuli): Use a watcher instead of polling.
		time.Sleep(10 * time.Second)

		// Make sure there's a pod.
		pod, err := kubeClient.CoreV1().Pods(workload.DefaultNamespace).Get(name, apismetav1.GetOptions{})
		if err != nil {
			return "", errors.WithStack(err)
		}

		// Make sure the pod isn't failing.
		if isFailing, err := workload.IsPodFailing(pod); isFailing {
			return "", err
		}

		if pod.Status.Phase == v1.PodRunning {
			return pod.Status.PodIP, nil
		}

		// return an error, if has not get the pod ip for 60 seconds.
		if n > 5 {
			return "", errors.Errorf("long time to get pod %s ip", name)
		}
		n++
	}
}

// findTemplate returns the true testing tool template and arguments for different testing cases.
func (t *TestingTool) findTemplate(name string) (string, string) {
	if t.CurrentTesting.Name == benchmarkPodIPDiffNode {
		return wrkPodAntiAffinity, t.CurrentTesting.TestingToolArgs
	}
	if t.CurrentTesting.Name == benchmarkPodIPSameNode {
		return wrkPodAffinity, t.CurrentTesting.TestingToolArgs
	}
	return "", ""
}

func (t *TestingTool) workloadPodName(testingName string) string {
	return strings.ToLower("capstan-" + t.Workload.GetName() + "-" + testingName)
}

func (t *TestingTool) testingPodName(testingName string) string {
	return strings.ToLower("capstan-" + t.GetName() + "-" + testingName)
}
