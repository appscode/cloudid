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

// https://developers.digitalocean.com/documentation/metadata/#metadata-in-json
func DetectDigitalOcean(done chan<- string) {
	md := struct {
		DropletID  int      `json:"droplet_id"`
		Hostname   string   `json:"hostname"`
		VendorData string   `json:"vendor_data"`
		PublicKeys []string `json:"public_keys"`
		Region     string   `json:"region"`
		Interfaces struct {
			Private []struct {
				Ipv4 struct {
					IPAddress string `json:"ip_address"`
					Netmask   string `json:"netmask"`
					Gateway   string `json:"gateway"`
				} `json:"ipv4"`
				Mac  string `json:"mac"`
				Type string `json:"type"`
			} `json:"private"`
			Public []struct {
				Ipv4 struct {
					IPAddress string `json:"ip_address"`
					Netmask   string `json:"netmask"`
					Gateway   string `json:"gateway"`
				} `json:"ipv4"`
				Ipv6 struct {
					IPAddress string `json:"ip_address"`
					Cidr      int    `json:"cidr"`
					Gateway   string `json:"gateway"`
				} `json:"ipv6"`
				Mac  string `json:"mac"`
				Type string `json:"type"`
			} `json:"public"`
		} `json:"interfaces"`
		FloatingIP struct {
			Ipv4 struct {
				Active bool `json:"active"`
			} `json:"ipv4"`
		} `json:"floating_ip"`
		DNS struct {
			Nameservers []string `json:"nameservers"`
		} `json:"dns"`
	}{}

	hc := httpclient.Default()
	resp, err := hc.Call(http.MethodGet, "http://169.254.169.254/metadata/v1.json", nil, &md, false)
	if err == nil && resp.StatusCode == http.StatusOK && md.DropletID > 0 {
		done <- "digitalOcean"
	}
	done <- ""
}
