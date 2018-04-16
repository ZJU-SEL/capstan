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
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/ZJU-SEL/capstan/pkg/capstan/types"
	"github.com/ZJU-SEL/capstan/pkg/util"
	"github.com/ZJU-SEL/capstan/pkg/workload"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// TestTool represents a test tool which used to test a workload managed by helm.
type TestTool struct {
	Workload    *Workload
	Name        string
	Script      string
	Steps       time.Duration
	CurrentTest workload.TestCase
	TestCaseSet []workload.TestCase
}

// Ensure the test tool implements workload.Tool interface.
var _ workload.Tool = &TestTool{}

// Run runs the defined test case set of the test tool (to adhere to workload.Tool interface).
func (t *TestTool) Run(kubeClient kubernetes.Interface, testCase workload.TestCase) error {
	t.CurrentTest = testCase

	// 1. create a namespace used run test
	err := workload.CreateNamespace(kubeClient, workload.Namespace)
	if err != nil {
		return errors.Wrap(err, "Failed create namespace")
	}

	// 2. start a workload for the test case.
	ret, err := util.RunCommand("helm", "install", "--name", t.Workload.Helm.Name, "--set", t.Workload.Helm.Set, "--namespace", workload.Namespace, t.Workload.Helm.Chart)
	if err != nil {
		return errors.Errorf("helm install failed, ret:%s, error:%v", strings.Join(ret, "\n"), err)
	}

	// 3. check workload is available or not.
	err = workload.CheckWorkloadAvailable(kubeClient, t)
	if err != nil {
		return errors.Wrapf(err, "workload %s is not available for testing case %s", t.Workload.Name, testCase.Name)
	}

	// 4. create ConfigMaps which is used by a test pod.
	// 4.1 create capstan-script configmap.
	err = workload.CreateConfigMapFromFile(kubeClient, t.Script)
	if err != nil {
		return errors.Wrapf(err, "failed create configmap from file %s", t.Script)
	}

	// 4.2 create capstan-envs configmap.
	data := parseEnvs(testCase.Envs)

	labelData := map[string]string{
		"job":          t.Workload.Name,
		"uid":          types.UUID,
		"provider":     types.Provider,
		"startTime":    time.Now().Format("2006-01-02 15:04:05"),
		"workloadName": t.Workload.Name,
		"testCase":     t.CurrentTest.Name,
		"affinity":     strconv.FormatBool(t.CurrentTest.Affinity),
		"metrics":      t.CurrentTest.Metrics,
	}
	data["PrometheusLabel"] = megeLabel(labelData)

	err = workload.CreateConfigMap(kubeClient, "capstan-envs", data)
	if err != nil {
		return errors.Wrapf(err, "failed create configmap from map %s", data)
	}

	// 5. start a test pod to test the workload.
	testPodName := workload.BuildTestPodName(t.GetName(), testCase.Name)
	testPod, args := t.findTemplate(testCase.Name)

	tempTestArgs := struct{ Name, Namespace, TestingName, Image, Label, Args string }{
		Name:        testPodName,
		Namespace:   workload.Namespace,
		TestingName: testCase.Name,
		Image:       "wadelee/capstan-base",
		Label:       t.Workload.Helm.Name + "-" + t.Workload.Name,
		Args:        workload.FomatArgs(args),
	}

	testPodBytes, err := workload.ParseTemplate(testPod, tempTestArgs)
	if err != nil {
		return errors.Wrapf(err, "unable to parse %v using %v", testPod, tempTestArgs)
	}

	glog.V(4).Infof("Creating testing pod %q of testing case %s", testPodName, testCase.Name)
	if err := workload.CreatePod(kubeClient, testPodBytes); err != nil {
		return errors.Wrapf(err, "unable to create the testing pod for testing case %s", testCase.Name)
	}

	return nil
}

// HasTestDone checks the test case has finished or not(use the finish mark "Capstan finish the test case"). (to adhere to workload.Tool interface).
func (t *TestTool) HasTestDone(kubeClient kubernetes.Interface) error {
	name := workload.BuildTestPodName(t.GetName(), t.CurrentTest.Name)
	for {
		// Sleep between each poll
		// TODO(mozhuli): Use a watcher instead of polling.
		time.Sleep(30 * time.Second)

		// Make sure there's a pod.
		pod, err := kubeClient.CoreV1().Pods(workload.Namespace).Get(name, apismetav1.GetOptions{})
		if err != nil {
			return errors.WithStack(err)
		}

		// Make sure the pod isn't failing.
		if isFailing, err := workload.IsPodFailing(pod); isFailing {
			return err
		}

		if pod.Status.Phase == v1.PodSucceeded {
			glog.V(4).Infof("The test case has finished successfully, starting fetch pod %s log", name)

			body, err := kubeClient.CoreV1().Pods(workload.Namespace).GetLogs(
				name,
				&v1.PodLogOptions{},
			).Do().Raw()

			if err != nil {
				return errors.Wrapf(err, "Failed fetch pod %s test log", name)
			}

			glog.V(5).Infof("The test pod %s's test log:\n%s", name, string(body))

			// export to capstan result directory.
			outdir := path.Join(types.ResultsDir, types.UUID, "workloads", t.Workload.Name, t.GetName(), t.CurrentTest.Name)
			if err = os.MkdirAll(outdir, 0755); err != nil {
				return errors.WithStack(err)
			}

			outfile := path.Join(outdir, t.GetName()) + ".log"
			if err = ioutil.WriteFile(outfile, body, 0644); err != nil {
				return errors.WithStack(err)
			}
			return nil
		} else if pod.Status.Phase == v1.PodFailed {
			return errors.Errorf("The test pod %s exited with a non-zero exit code or was stopped by the system. Please check your test script %q is correct ?", name, t.Script)
		}
	}
}

// Cleanup cleans up all resources created by a test case for mysql test tool (to adhere to workload.Tool interface).
func (t *TestTool) Cleanup(kubeClient kubernetes.Interface) error {
	// Delete the release of the workload
	ret, err := util.RunCommand("helm", "delete", "--purge", t.Workload.Helm.Name)
	if err != nil {
		return errors.Errorf("helm install failed, ret:%s, error:%v", strings.Join(ret, "\n"), err)
	}

	// Delete namespace to clean resouces quickly
	err = workload.DeleteNamespace(kubeClient, workload.Namespace)
	if err != nil {
		return errors.Wrapf(err, "Failed delete namespace %s", workload.Namespace)
	}
	return nil
}

// GetName returns the name of the testing tool (to adhere to workload.Tool interface).
func (t *TestTool) GetName() string {
	return t.Name
}

// GetSteps returns the steps between each testing case (to adhere to workload.Tool interface).
func (t *TestTool) GetSteps() time.Duration {
	return t.Steps
}

// GetTestCaseSet returns the test case set which the test tool will run (to adhere to workload.Tool interface).
func (t *TestTool) GetTestCaseSet() []workload.TestCase {
	return t.TestCaseSet
}

// GetWorkload returns a workload that the test tool to be test (to adhere to workload.Tool interface).
func (t *TestTool) GetWorkload() workload.Workload {
	return t.Workload.workload
}

// findTemplate returns the true testing tool template and arguments for different testing cases.
func (t *TestTool) findTemplate(name string) (string, string) {
	if t.CurrentTest.Affinity {
		return PodAffinity, t.CurrentTest.Args
	}
	return PodAntiAffinity, t.CurrentTest.Args
}

// parseEnvs parse string to map[string]string
func parseEnvs(raw string) map[string]string {
	data := make(map[string]string)
	if raw != "" {
		for _, re := range strings.Split(raw, ",") {
			env := strings.Split(re, "=")
			data[env[0]] = env[1]
		}
	}
	data["PushgatewayEndpoint"] = types.PushgatewayEndpoint
	return data
}

func megeLabel(data map[string]string) string {
	var str string
	for k, v := range data {
		str += k + "=" + v + ","
	}

	return strings.TrimSuffix(str, ",")
}
