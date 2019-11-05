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
	"fmt"
	"net"
	"net/http"
	"time"
)

// https://github.com/scaleway/initrd/issues/84
func DetectScaleway(done chan<- string) {
	for port := 1; port <= 1024; port++ {
		hc := &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyFromEnvironment,
				Dial: (&net.Dialer{
					LocalAddr: &net.TCPAddr{
						Port: port,
					},
					Timeout:   5 * time.Second,
					KeepAlive: 5 * time.Second,
				}).Dial,
				TLSHandshakeTimeout: 5 * time.Second,
			},
		}
		resp, err := hc.Get("http://169.254.42.42/user_data")
		if err != nil {
			fmt.Printf("Bind to local port %d failed: %s\n", port, err.Error())
		} else {
			if resp.StatusCode == http.StatusOK {
				done <- "scaleway"
			} else {
				done <- ""
			}
			return
		}
	}
	done <- ""
}
