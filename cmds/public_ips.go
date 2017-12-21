package cmds

import (
	"fmt"
	"os"
	"strings"

	"github.com/appscode/go/net"
	"github.com/spf13/cobra"
)

func NewCmdPublicIPs() *cobra.Command {
	all := true
	routable := false
	cmd := &cobra.Command{
		Use:               "public-ips",
		Short:             "Prints public ip(s) for current host",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			var ips []string
			var err error
			if routable {
				ips, _, err = net.RoutableIPs()
				if err != nil {
					Fatal(fmt.Errorf("failed to detect routable ips. Reason: %v", err))
				}
			} else {
				ips, _, err = net.HostIPs()
				if err != nil {
					Fatal(fmt.Errorf("failed to detect host ips. Reason: %v", err))
				}
			}
			if !all && len(ips) > 0 {
				fmt.Print(ips[0])
			} else {
				fmt.Print(strings.Join(ips, ","))
			}
			os.Exit(0)
		},
	}
	cmd.Flags().BoolVar(&all, "all", all, "If true, returns all private ips.")
	cmd.Flags().BoolVar(&routable, "routable", routable, "If true, uses https://ipinfo.io/ip to detect ip in case host interfaces are missing public ips.")
	return cmd
}
