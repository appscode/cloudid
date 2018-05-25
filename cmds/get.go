package cmds

import (
	onessl "github.com/appscode/onessl/cmds"
	"github.com/spf13/cobra"
)

func NewCmdGet() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "get",
		Short:             `Get stuff`,
		DisableAutoGenTag: true,
	}
	cmd.AddCommand(onessl.NewCmdGetCACert())
	cmd.AddCommand(NewCmdGetPublicKey())
	return cmd
}
