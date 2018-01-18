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

package main

import (
	"fmt"
	"os"

	"github.com/ZJU-SEL/capstan/pkg/capstan"
	"github.com/ZJU-SEL/capstan/pkg/util"
	"github.com/golang/glog"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
)

var (
	kubeconfig    = pflag.String("kubeconfig", "/etc/kubernetes/admin.conf", "path to kubernetes admin config file")
	capstanConfig = pflag.String("capstanconfig", "/etc/capstan.conf", "path to capstan config file")
	version       = pflag.Bool("version", false, "Display version")
	// VERSION is the version of capstan.
	VERSION = "1.0"
)

func initK8sClient() (*kubernetes.Clientset, error) {
	// Create kubernetes client config. Use kubeconfig if given, otherwise assume in-cluster.
	config, err := util.NewClusterConfig(*kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build kubeconfig: %v", err)
	}
	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes clientset: %v", err)
	}

	return kubeClient, nil
}

func main() {
	util.InitFlags()
	util.InitLogs()
	defer util.FlushLogs()

	// Print capstan version
	if *version {
		fmt.Println(VERSION)
		os.Exit(0)
	}

	// Initilize kubernetes client
	kubeClient, err := initK8sClient()
	if err != nil {
		glog.Fatal(err)
	}

	// Run capstan
	if err := capstan.Run(kubeClient, *capstanConfig); err != nil {
		glog.Fatal(err)
	}
}
