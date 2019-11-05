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
package cmds

import (
	"fmt"
	"os"
	"strings"

	"github.com/appscode/go/net"
	"github.com/spf13/cobra"
)

func NewCmdLinodeHostname() *cobra.Command {
	var (
		clusterName string
	)
	cmd := &cobra.Command{
		Use:               "hostname",
		Short:             "Prints hostname based on public IP for current Linode host",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			ips, _, err := net.HostIPs()
			if err != nil {
				Fatal(fmt.Errorf("failed to detect host ips. Reason: %v", err))
			}
			if len(ips) == 0 {
				os.Exit(1)
			}
			parts := strings.SplitN(ips[0], ".", 4)
			fmt.Printf("%s-%03s-%03s-%03s-%03s", clusterName, parts[0], parts[1], parts[2], parts[3])
			os.Exit(0)
		},
	}
	cmd.Flags().StringVarP(&clusterName, "cluster", "k", "", "Name of cluster")
	return cmd
}
