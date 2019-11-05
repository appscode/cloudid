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
package lib

import (
	"os/exec"

	"pharmer.dev/pre-k/internal"
)

var ExecCommand = exec.Command

func DetectCloudProvider() string {
	done := make(chan string)
	go internal.DetectAWS(done)
	go internal.DetectGCE(done)
	go internal.DetectDigitalOcean(done)
	go internal.DetectAzure(done)
	go internal.DetectVultr(done)
	go internal.DetectLinode(done)
	go internal.DetectSoftlayer(done)
	go internal.DetectScaleway(done)

	n := 8 // total number of go routines
	i := 0
	provider := ""
	for ; i < n; i++ {
		p := <-done
		if p != provider {
			provider = p
			break
		}
	}
	if i < n {
		// run drainer
		go func() {
			for ; i < n; i++ {
				<-done
			}
		}()
	}
	return provider
}

func UserExists(u string) bool {
	return ExecCommand("id", "-u", u).Run() == nil
}

func Run(cmd string, args ...string) error {
	return ExecCommand(cmd, args...).Run()
}
