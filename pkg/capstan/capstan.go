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
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ZJU-SEL/capstan/pkg/capstan/loader"
	"github.com/ZJU-SEL/capstan/pkg/capstan/types"
	"github.com/ZJU-SEL/capstan/pkg/dashboard"
	"github.com/golang/glog"
	"github.com/pkg/errors"
	"k8s.io/client-go/kubernetes"
)

// Run is the entry to run capstan.
//
// Basic workflow:
//
// 1. Read capstan config
// 2. Load all workloads
// 3. Start runs all testing workloads sequentially
// 4. Launch the HTTP server
func Run(kubeClient kubernetes.Interface, capstanConfig string) error {
	// 1. Read capstan config.
	cfg, err := types.ReadConfig(capstanConfig)
	if err != nil {
		return errors.Wrap(err, "Failed read capstan config")
	}
	glog.V(1).Infof("Initializing capstan with config %v", cfg)

	if len(cfg.Workloads) == 0 {
		return errors.New("Testing workload not set, exit")
	}

	// 2. Load all workloads
	workloads, err := loader.LoadAllWorkloads(cfg.Workloads)
	if err != nil {
		return errors.Wrap(err, "Failed load workloads")
	}

	// 3. Start runs all testing workloads sequentially
	testingDone := make(chan bool)
	testingErr := make(chan error)
	go func() {
		for _, wk := range workloads {
			err := wk.Run(kubeClient)
			if err != nil {
				testingErr <- err
			}
			time.Sleep(time.Duration(cfg.Steps) * time.Second)
		}
		testingDone <- true
	}()

	// 4. Launch the HTTP server
	srv := &http.Server{
		Addr:    cfg.Address,
		Handler: dashboard.NewHandler(),
	}
	doneServ := make(chan error)
	go func() {
		doneServ <- srv.ListenAndServe()
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-testingDone:
		glog.V(4).Info("Finished all tests")
	case <-term:
		glog.V(4).Info("Received SIGTERM, exiting gracefully...")
		return deleteAllResources(kubeClient)
	case err := <-testingErr:
		return err
	case err := <-doneServ:
		return err
	}
	return nil
}

func deleteAllResources(kubeClient kubernetes.Interface) error {
	return errors.Errorf("Not Implemented")
}
