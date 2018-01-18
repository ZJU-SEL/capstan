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

package dashboard

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// Handler is a net/http Handler that can handle API requests of dashboard.
type Handler struct {
	mux.Router
}

// NewHandler register all the dashboard apis.
func NewHandler() http.Handler {
	handler := &Handler{
		Router: *mux.NewRouter(),
	}
	handler.HandleFunc("/overview", overviewHandler).Methods("GET")
	handler.HandleFunc("/download", downloadHandler).Methods("GET")
	return handler
}

// overviewHandler hands the overview request.
func overviewHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Not Implemented")
}

// downloadHandler hands the download request.
func downloadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Not Implemented")
}
