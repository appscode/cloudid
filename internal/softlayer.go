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

// https://sldn.softlayer.com/blog/jarteche/getting-started-user-data-and-post-provisioning-scripts
// https://github.com/bodenr/cci/wiki/SL-user-metadata
func DetectSoftlayer(done chan<- string) {
	hc := httpclient.Default()
	resp, err := hc.Call(http.MethodGet, "https://api.service.softlayer.com/rest/v3/SoftLayer_Resource_Metadata/UserMetadata.txt", nil, nil, false)
	if err == nil &&
		(resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound) {
		done <- "softlayer"
	}
	done <- ""
}
