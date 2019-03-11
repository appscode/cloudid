package cmds

import (
	"flag"
	"log"
	"strings"

	v "github.com/appscode/go/version"
	ga "github.com/jpillora/go-ogle-analytics"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"kmodules.xyz/client-go/logs"
)

const (
	gaTrackingCode = "UA-62096468-20"
)

func NewRootCmd(version string) *cobra.Command {
	var (
		enableAnalytics = true
	)
	rootCmd := &cobra.Command{
		Use:               "pre-k [command]",
		Short:             `Pre-k by AppsCode - Utilities for your cloud`,
		DisableAutoGenTag: true,
		PersistentPreRun: func(c *cobra.Command, args []string) {
			c.Flags().VisitAll(func(flag *pflag.Flag) {
				log.Printf("FLAG: --%s=%q", flag.Name, flag.Value)
			})
			if enableAnalytics && gaTrackingCode != "" {
				if client, err := ga.NewClient(gaTrackingCode); err == nil {
					parts := strings.Split(c.CommandPath(), " ")
					client.Send(ga.NewEvent(parts[0], strings.Join(parts[1:], "/")).Label(version))
				}
			}
		},
	}
	rootCmd.PersistentFlags().BoolVar(&enableAnalytics, "analytics", enableAnalytics, "Send analytical events to Google Guard")
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	logs.ParseFlags()

	rootCmd.AddCommand(NewCmdCreate())
	rootCmd.AddCommand(NewCmdCheck())
	rootCmd.AddCommand(NewCmdGet())
	rootCmd.AddCommand(NewCmdLinode())
	rootCmd.AddCommand(NewCmdMachine())
	rootCmd.AddCommand(NewCmdMerge())
	rootCmd.AddCommand(NewCmdVultr())
	rootCmd.AddCommand(v.NewCmdVersion())
	rootCmd.AddCommand(NewCmdMountMasterPD())
	return rootCmd
}
