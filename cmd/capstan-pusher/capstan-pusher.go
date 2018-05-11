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
	"flag"
	"fmt"
	"os"

	"github.com/ZJU-SEL/capstan/pkg/capstan-pusher"
)

func main() {
	PGWEndpoint := flag.String("endpoint", "", "Path to the pushGateWayls.")
	flag.Parse()

	if *PGWEndpoint == "" {
		fmt.Fprintf(os.Stderr, "%s\n", "No phase named endpoint found")
		os.Exit(1)
	}

	if len(flag.Args()) != 1 {
		fmt.Fprintf(os.Stderr, "%s\n", "The format of the result cannot be resolved.")
		os.Exit(1)
	}

	err := push.Push(flag.Args()[0], *PGWEndpoint)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
