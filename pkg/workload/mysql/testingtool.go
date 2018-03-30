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

package mysql

import (
	"bufio"
	"bytes"
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
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/push"
	v1 "k8s.io/api/core/v1"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	toolName              = "tpcc-mysql"
	benchmarkTPMCSameNode = "benchmarkTPMCSameNode"
	benchmarkTPMCDiffNode = "benchmarkTPMCDiffNode"
)

// TestingCaseSet is the list of mysql defined testing cases.
var TestingCaseSet = []string{
	"benchmarkTPMCSameNode",
	"benchmarkTPMCDiffNode",
}

// TestingTool represents the mysql testing tool.
type TestingTool struct {
	Workload       *Workload
	Name           string
	Image          string
	Steps          time.Duration
	StartTime      time.Time
	CurrentTesting workload.TestingCase
	TestingCaseSet []workload.TestingCase
}

// Ensure mysql testing tool implements workload.Tool interface.
var _ workload.Tool = &TestingTool{}

// Run runs the defined testing case set for mysql testing tool (to adhere to workload.Tool interface).
func (t *TestingTool) Run(kubeClient kubernetes.Interface, testingCase workload.TestingCase) error {
	t.CurrentTesting = testingCase
	t.StartTime = time.Now()

	// 1. start a workload for the testing case.
	ret, err := util.RunCommand("helm", "install", "--name", t.Workload.Helm.Name, "--set", t.Workload.Helm.Set, "--namespace", workload.Namespace, t.Workload.Helm.Chart)
	if err != nil {
		return errors.Errorf("helm install failed, ret:%s, error:%v", strings.Join(ret, "\n"), err)
	}

	// 2. check the deployment of the workload is available.
	glog.V(4).Infof("Check the deployment of %s workload is available or not", t.Workload.Name)
	err = workload.CheckDeployment(kubeClient, t.Workload.Helm.Name+"-mysql")
	if err != nil {
		return errors.Wrapf(err, "unable to check the deployment of %s workload is available or not for testing case %s", t.Workload.Name, testingCase.Name)
	}

	// 3. start a testing pod for testing the workload.
	testingPodName := workload.BuildTestingPodName(t.GetName(), testingCase.Name)
	testingPod, args := t.findTemplate(testingCase.Name)
	password, err := getMysqlRootPassword(t.Workload.Helm.Set)
	if err != nil {
		return err
	}
	tempTestingArgs := struct{ Name, Namespace, TestingName, Image, Label, Args, DNSName, PASSWORD string }{
		Name:        testingPodName,
		Namespace:   workload.Namespace,
		TestingName: testingCase.Name,
		Image:       t.GetImage(),
		Label:       t.Workload.Helm.Name + "-mysql",
		Args:        workload.FomatArgs(args),
		DNSName:     t.Workload.Helm.Name + "-mysql",
		PASSWORD:    password,
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

// GetTestingResults gets the testing results of mysql testing case (to adhere to workload.Tool interface).
func (t *TestingTool) GetTestingResults(kubeClient kubernetes.Interface) error {
	name := workload.BuildTestingPodName(t.GetName(), t.CurrentTesting.Name)
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

		// Check testing has done.
		body, err := kubeClient.CoreV1().Pods(workload.Namespace).GetLogs(
			name,
			&v1.PodLogOptions{},
		).Do().Raw()

		if err != nil {
			return errors.WithStack(err)
		}

		glog.V(5).Infof("Checking testing has done:\n%s", string(body))
		if workload.HasTestingDone(body) {
			glog.V(4).Infof("Testing case %s has done", t.CurrentTesting.Name)

			// export to capstan result directory.
			outdir := path.Join(types.ResultsDir, types.UUID, "workloads", t.Workload.Name, t.GetName(), t.CurrentTesting.Name)
			if err = os.MkdirAll(outdir, 0755); err != nil {
				return errors.WithStack(err)
			}

			outfile := path.Join(outdir, t.GetName()) + ".log"
			if err = ioutil.WriteFile(outfile, body, 0644); err != nil {
				return errors.WithStack(err)
			}

			// export to prometheus pushGateway.
			data, err := getTPMC(body)
			if err != nil {
				return errors.Wrapf(err, "Failed to get tpmc")
			}

			tpmc := prometheus.NewGauge(prometheus.GaugeOpts{
				Name: "capstan_mysql_tpmc",
				Help: "The tpmc of mysql testing case",
			})
			tpmc.Set(data)
			if err := push.Collectors(
				"mysql",
				map[string]string{
					"uid":          types.UUID,
					"provider":     types.Provider,
					"startTime":    t.StartTime.Format("2006-01-02 15:04:05"),
					"endTime":      time.Now().Format("2006-01-02 15:04:05"),
					"testingNode":  pod.Status.HostIP,
					"workloadName": t.Workload.Name,
					"testingName":  t.GetName(),
					"testingCase":  t.CurrentTesting.Name,
				},
				types.PushgatewayEndpoint,
				tpmc,
			); err != nil {
				return errors.Wrapf(err, "Could not push metrics to Pushgateway")
			}

			return nil
		}
	}
}

// Cleanup cleans up all resources created by a testing case for mysql testing tool (to adhere to workload.Tool interface).
func (t *TestingTool) Cleanup(kubeClient kubernetes.Interface) error {
	// Delete the release of the workload
	ret, err := util.RunCommand("helm", "delete", "--purge", t.Workload.Helm.Name)
	if err != nil {
		return errors.Errorf("helm install failed, ret:%s, error:%v", strings.Join(ret, "\n"), err)
	}
	// Delete testing pod.
	if err := workload.DeletePod(kubeClient, workload.BuildTestingPodName(t.GetName(), t.CurrentTesting.Name)); err != nil {
		return err
	}
	return nil
}

// GetName returns the name of mysql testing tool (to adhere to workload.Tool interface).
func (t *TestingTool) GetName() string {
	return t.Name
}

// GetImage returns the image name of mysql testing tool (to adhere to workload.Tool interface).
func (t *TestingTool) GetImage() string {
	return t.Image
}

// GetSteps returns the steps between each testing case (to adhere to workload.Tool interface).
func (t *TestingTool) GetSteps() time.Duration {
	return t.Steps
}

// GetTestingCaseSet returns the testing case set which the mysql testing tool will run (to adhere to workload.Tool interface).
func (t *TestingTool) GetTestingCaseSet() []workload.TestingCase {
	return t.TestingCaseSet
}

// findTemplate returns the true testing tool template and arguments for different testing cases.
func (t *TestingTool) findTemplate(name string) (string, string) {
	if t.CurrentTesting.Name == benchmarkTPMCDiffNode {
		return mysqlTPCCPodAntiAffinity, t.CurrentTesting.TestingToolArgs
	}
	if t.CurrentTesting.Name == benchmarkTPMCSameNode {
		return mysqlTPCCPodAffinity, t.CurrentTesting.TestingToolArgs
	}
	return "", ""
}

func getTPMC(data []byte) (float64, error) {
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, " TpmC") {
			tpmc, err := strconv.ParseFloat(strings.Fields(line)[0], 64)
			if err != nil {
				return 0, errors.WithStack(err)
			}
			return tpmc, nil
		}
	}
	return 0, errors.Errorf("results not contain TpmC")
}

func getMysqlRootPassword(set string) (string, error) {
	if !strings.Contains(set, "mysqlRootPassword") {
		return "", errors.Errorf("helm'set section should contain %q", "mysqlRootPassword")
	}
	for _, option := range strings.Split(set, ",") {
		if strings.Contains(option, "mysqlRootPassword") {
			return strings.Split(option, "=")[1], nil
		}
	}
	return "", errors.Errorf("helm'set section should contain %q", "mysqlRootPassword")
}
