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

package capstan

import (
	"net/http"
	"time"

	"github.com/ZJU-SEL/capstan/pkg/analysis"
	"github.com/ZJU-SEL/capstan/pkg/capstan/loader"
	"github.com/ZJU-SEL/capstan/pkg/capstan/types"
	"github.com/ZJU-SEL/capstan/pkg/dashboard"
	"github.com/ZJU-SEL/capstan/pkg/data/cadvisor"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

// Run is the entry to run capstan.
//
// Basic workflow:
//
// 1. Load all workloads
// 2. Start obtaining cadvisor data
// 3. Start obtaining resource usage data of kubelet
// 4. Start runs testing workload sequentially
// 5. Launch the HTTP server
// 6. Block analysis until a testing workload results has been returned
func Run(kubeClient kubernetes.Interface, capstanConfig string) error {
	// Read capstan config.
	cfg, err := types.ReadConfig(capstanConfig)
	if err != nil {
		return errors.Wrap(err, "Failed read capstan config")
	}
	glog.V(1).Infof("Initializing capstan with config %v", cfg)

	if len(cfg.Workloads) == 0 {
		return errors.New("Testing workload not set, exit")
	}

	// 1. Load all workloads
	workloads, err := loader.LoadAllWorkloads(cfg.Workloads)
	if err != nil {
		return errors.Wrap(err, "Failed load workloads")
	}

	// 2. Start obtaining cadvisor data
	cadvisorErr := make(chan error)
	glog.V(1).Infof("Starting obtaining cadvisor data")
	go func() {
		cadvisorErr <- cadvisor.Start(cfg.Cadvisor)
	}()

	// 3. TODO(mozhuli): Start obtaining resource usage data of kubelet

	// 4. Start runs all testing workloads sequentially
	doneTesting := make(chan bool, 1)
	testingErr := make(chan error)
	go func() {
		for _, wk := range workloads {
			err := wk.Run(kubeClient)
			if err != nil {
				testingErr <- err
			}
			time.Sleep(time.Duration(cfg.Steps) * time.Second)
		}
		doneTesting <- true
	}()

	// 5. Launch the HTTP server
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: dashboard.NewHandler(),
	}
	doneServ := make(chan error)
	go func() {
		doneServ <- srv.ListenAndServe()
	}()

	// 6. Block analysis until a testing workload results has been returned
	analysisErr := make(chan error)
	go func() {
		if <-doneTesting {
			analysisErr <- analysis.Start(cfg.ResultsDir)
		}
	}()

	select {
	case err := <-cadvisorErr:
		return err
	case err := <-testingErr:
		return err
	case err := <-analysisErr:
		return err
	case err := <-doneServ:
		if err != nil {
			return err
		}
	}
	return nil
}
