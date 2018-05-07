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
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/ZJU-SEL/capstan/pkg/util"
	"github.com/pkg/errors"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	apismetav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kuberuntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
)

// ParseTemplate uses the obj to parse the strtmpl template.
func ParseTemplate(strtmpl string, obj interface{}) ([]byte, error) {
	var buf bytes.Buffer
	tmpl, err := template.New("template").Parse(strtmpl)
	if err != nil {
		return nil, errors.Wrap(err, "error when parsing template")
	}
	err = tmpl.Execute(&buf, obj)
	if err != nil {
		return nil, errors.Wrap(err, "error when executing template")
	}
	return buf.Bytes(), nil
}

// CreatePod creates a pod using podBytes.
func CreatePod(kubeClient kubernetes.Interface, podBytes []byte) error {
	pod := &v1.Pod{}
	if err := kuberuntime.DecodeInto(scheme.Codecs.UniversalDecoder(), podBytes, pod); err != nil {
		return errors.Wrap(err, "unable to decode pod")
	}

	_, err := kubeClient.CoreV1().Pods(Namespace).Create(pod)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// DeletePod deletes a pod with the name.
func DeletePod(kubeClient kubernetes.Interface, name string) error {
	if err := kubeClient.CoreV1().Pods(Namespace).Delete(name, apismetav1.NewDeleteOptions(0)); err != nil {
		return errors.Wrapf(err, "failed to delete pod %v", name)
	}

	err := wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		_, err := kubeClient.CoreV1().Pods(Namespace).Get(name, apismetav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return true, nil
			}

			return false, err
		}

		return false, nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// IsPodFailing returns whether a test case pod is failing and isn't likely to succeed.
// TODO(mozhuli): this may require more revisions as we get more experience with
// various types of failures that can occur.
func IsPodFailing(pod *v1.Pod) (bool, error) {
	// Check if the pod is unschedulable
	for _, cond := range pod.Status.Conditions {
		if cond.Reason == "Unschedulable" {
			return true, errors.Errorf("Can't schedule pod: %v", cond.Message)
		}
	}

	for _, cstatus := range pod.Status.ContainerStatuses {
		// Check if a container in the pod is restarting multiple times
		if cstatus.RestartCount > 2 {
			return true, errors.Errorf("Container %v has restarted unsuccessfully %v times", cstatus.Name, cstatus.RestartCount)

		}

		// Check if it can't fetch its image
		if waiting := cstatus.State.Waiting; waiting != nil {
			if waiting.Reason == "ImagePullBackOff" || waiting.Reason == "ErrImagePull" {
				return true, errors.Errorf("Container %v is in state %v", cstatus.Name, waiting.Reason)

			}
		}
	}

	return false, nil
}

// CheckWorkloadAvailable check workload is available or not.
// TODO(mozhuli): add more rules to check workload is available or not.
func CheckWorkloadAvailable(kubeClient kubernetes.Interface, tool Tool) error {
	switch tool.GetWorkload().Name {
	case "nginx":
		return checkDeployment(kubeClient, tool.GetWorkload().Helm.Name+"-"+tool.GetWorkload().Name)
	case "mysql":
		return checkDeployment(kubeClient, tool.GetWorkload().Helm.Name+"-"+tool.GetWorkload().Name)
	case "iperf3":
		return checkDeployment(kubeClient, tool.GetWorkload().Helm.Name+"-"+tool.GetWorkload().Name)
	case "spark":
		return checkDeployment(kubeClient, tool.GetWorkload().Helm.Name+"-"+"master")
	case "kubeflow":
		err := checkDeployment(kubeClient, tool.GetWorkload().Helm.Name+"-"+tool.GetWorkload().Name+"-"+"tf-job-operator")
		if err != nil {
			return err
		}
		err = checkDeployment(kubeClient, tool.GetWorkload().Helm.Name+"-"+tool.GetWorkload().Name+"-"+"spartakus-volunteer")
		if err != nil {
			return err
		}
		return checkDeployment(kubeClient, tool.GetWorkload().Helm.Name+"-"+tool.GetWorkload().Name+"-"+"ambassador")
	}
	return errors.Errorf("Not meet any rules to check the workload available or not")
}

// checkDeployment check the deployment created by a workload available or not.
func checkDeployment(kubeClient kubernetes.Interface, name string) error {
	n := 0
	for {
		// Sleep between each poll, which should give the workload enough time to create
		// TODO(mozhuli): Use a watcher instead of polling.
		time.Sleep(30 * time.Second)

		// Make sure there's a deployment.
		deployment, err := kubeClient.AppsV1().Deployments(Namespace).Get(name, apismetav1.GetOptions{})
		if err != nil {
			return errors.WithStack(err)
		}

		// Make sure the deployment is available.
		if deployment.Status.Conditions[0].Type == appsv1.DeploymentAvailable {
			return nil
		}

		// return an error, if deployment not available for 60 seconds.
		if n > 5 {
			return errors.Errorf("deployment %q not available for 60 seconds", name)
		}
		n++
	}
}

// FomatArgs fomats the config agrs to yaml agrs of kubernetes.
func FomatArgs(agrs string) string {
	ss := strings.Split(agrs, " ")
	var str string
	for i, s := range ss {
		if i == len(ss)-1 {
			str = str + fmt.Sprintf("\"%s\"", s)
		} else {
			str = str + fmt.Sprintf("\"%s\",", s)
		}
	}
	return str
}

// BuildTestPodName builds the name of test pod.
func BuildTestPodName(name, testName string) string {
	return strings.ToLower("capstan-" + name + "-" + testName)
}

// CreateNamespace creates a namespace.
func CreateNamespace(kubeClient kubernetes.Interface, namespace string) error {
	nsSpec := &v1.Namespace{ObjectMeta: apismetav1.ObjectMeta{Name: namespace}}
	_, err := kubeClient.CoreV1().Namespaces().Create(nsSpec)
	if err != nil {
		return errors.Wrapf(err, "Failed to create namespace %v", namespace)
	}

	return nil
}

// DeleteNamespace deletes a namespace.
func DeleteNamespace(kubeClient kubernetes.Interface, namespace string) error {
	if err := kubeClient.CoreV1().Namespaces().Delete(namespace, apismetav1.NewDeleteOptions(0)); err != nil {
		return errors.Wrapf(err, "Failed to delete namespace %s", namespace)
	}

	err := wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		_, err := kubeClient.CoreV1().Namespaces().Get(Namespace, apismetav1.GetOptions{})
		if err != nil {
			if apierrors.IsNotFound(err) {
				return true, nil
			}

			return false, err
		}

		return false, nil
	})
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// CleanNamespace make sure there is a clean namespace.
func CleanNamespace(kubeClient kubernetes.Interface, namespace string) error {
	_, err := kubeClient.CoreV1().Namespaces().Get(Namespace, apismetav1.GetOptions{})
	if err != nil {
		if apierrors.IsNotFound(err) {
			return nil
		}
	}

	return DeleteNamespace(kubeClient, namespace)
}

// CreateConfigMapFromFile creates a configmap from specific file.
func CreateConfigMapFromFile(kubeClient kubernetes.Interface, filePath string) error {
	ret, err := util.RunCommand("kubectl", "create", "configmap", "capstan-script", "--from-file=run_test.sh="+filePath, "-n="+Namespace)
	if err != nil {
		return errors.Errorf("failed create configmap, ret:%s, error:%v", strings.Join(ret, "\n"), err)
	}

	return nil
}

// CreateConfigMap creates a configmap from a map.
func CreateConfigMap(kubeClient kubernetes.Interface, name string, data map[string]string) error {
	cmSpec := &v1.ConfigMap{ObjectMeta: apismetav1.ObjectMeta{Name: name}, Data: data}
	_, err := kubeClient.CoreV1().ConfigMaps(Namespace).Create(cmSpec)
	if err != nil {
		return errors.Wrapf(err, "Failed to create configmap %v", name)
	}

	return nil
}
