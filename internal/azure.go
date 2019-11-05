/*
Copyright The Pharmer Authors.

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
package internal

import (
	"net/http"

	"github.com/appscode/go/net/httpclient"
)

// https://azure.microsoft.com/en-us/blog/what-just-happened-to-my-vm-in-vm-metadata-service/
func DetectAzure(done chan<- string) {
	md := struct {
		ID string `json:"ID"`
		UD string `json:"UD"`
		FD string `json:"FD"`
	}{}

	hc := httpclient.Default()
	resp, err := hc.Call(http.MethodGet, "http://169.254.169.254/metadata/v1/InstanceInfo", nil, &md, false)
	if err == nil && resp.StatusCode == http.StatusOK && md.ID != "" {
		done <- "azure"
	}
	done <- ""
}
