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

	"github.com/appscode/go/net"
	"github.com/spf13/cobra"
)

func NewCmdNodeIP() *cobra.Command {
	var ifaces []string
	cmd := &cobra.Command{
		Use:   "node-ip",
		Short: `Prints a IPv4 address for current host`,
		Long: `Prints a IPv4 address for current host for a given set of interface names. It always prefers a private IP over a public IP.
If no interface name is given, all interfaces are checked.`,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			_, ip, err := net.NodeIP(ifaces...)
			if err != nil {
				Fatal(fmt.Errorf("failed to detect node ip. Reason: %v", err))
			}
			fmt.Print(ip.String())
			os.Exit(0)
		},
	}
	cmd.Flags().StringSliceVar(&ifaces, "ifaces", ifaces, "Comma separated list of interface names. If empty, all interfaces are checked.")
	return cmd
}
