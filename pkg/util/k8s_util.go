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

package util

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	defaultQPS   = 100
	defaultBurst = 100
)

// NewClusterConfig builds a kubernetes cluster config.
func NewClusterConfig(kubeConfig string) (*rest.Config, error) {
	var cfg *rest.Config
	var err error

	if kubeConfig != "" {
		cfg, err = clientcmd.BuildConfigFromFlags("", kubeConfig)
	} else {
		cfg, err = rest.InClusterConfig()
	}
	if err != nil {
		return nil, err
	}

	// Setup default QPS and burst.
	cfg.QPS = defaultQPS
	cfg.Burst = defaultBurst
	return cfg, nil
}
