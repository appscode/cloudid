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
	"net"
	"strings"

	n2 "github.com/appscode/go/net"
)

func DetectLinode(done chan<- string) {
	externalIPs, _, err := n2.HostIPs()
	if err != nil {
		done <- ""
		return
	}
	for _, ip := range externalIPs {
		names, err := net.LookupAddr(ip)
		if err == nil {
			for _, name := range names {
				if strings.HasSuffix(name, ".members.linode.com.") {
					done <- "linode"
					return
				}
			}
		}
	}
	done <- ""
}
