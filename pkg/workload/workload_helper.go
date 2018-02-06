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
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"time"

	"github.com/pkg/errors"
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

	_, err := kubeClient.CoreV1().Pods(DefaultNamespace).Create(pod)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// DeletePod deletes a pod with the name.
func DeletePod(kubeClient kubernetes.Interface, name string) error {
	if err := kubeClient.CoreV1().Pods(DefaultNamespace).Delete(name, apismetav1.NewDeleteOptions(0)); err != nil {
		return errors.Wrapf(err, "failed to delete pod %v", name)
	}

	err := wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		_, err := kubeClient.CoreV1().Pods(DefaultNamespace).Get(name, apismetav1.GetOptions{})
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

// IsPodFailing returns whether a testing case pod is failing and isn't likely to succeed.
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

// HasTestingDone checks the testing case has finished
// or not(use the finish mark "Capstan Testing Done").
func HasTestingDone(data []byte) bool {
	scanner := bufio.NewScanner(bytes.NewBuffer(data))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "Capstan Testing Done" {
			return true
		}
	}
	return false
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
