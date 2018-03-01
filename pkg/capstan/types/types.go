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

package types

import (
	"encoding/json"
	"io/ioutil"

	"github.com/ZJU-SEL/capstan/pkg/prometheus"
	"github.com/ZJU-SEL/capstan/pkg/workload"
	"github.com/pkg/errors"
	"github.com/satori/go.uuid"
)

var (
	// ResultsDir is the directory of testing results.
	ResultsDir = "/tmp/capstan"
	// PushgatewayEndpoint is the endpoint of pushGateway.
	PushgatewayEndpoint string
	// Namespace is the namespace of capstan.
	Namespace = "capstan"
	// UUID is used to mark a run of capstan.
	UUID string
)

// Config is the internal representation of capstan configuration.
type Config struct {
	UUID       string `json:"UUID"`
	ResultsDir string `json:"ResultsDir"`
	Address    string `json:"Address"`
	Steps      int    `json:"Steps"`
	Namespace  string `json:"Namespace"`
	Prometheus prometheus.Config
	Workloads  []workload.Workload
}

// ReadConfig reads from a file with the given name and returns
// a config or an error if the file was unable to be parsed.
func ReadConfig(filepath string) (Config, error) {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return Config{}, errors.WithStack(err)
	}
	config := Config{}
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, errors.WithStack(err)
	}

	if config.UUID != "" {
		UUID = config.UUID
	} else {
		UUID = uuid.NewV4().String()
	}

	if config.ResultsDir != "" {
		ResultsDir = config.ResultsDir
	}

	if config.Namespace != "" {
		Namespace = config.Namespace
	}

	PushgatewayEndpoint = config.Prometheus.PushgatewayEndpoint

	return config, nil
}
