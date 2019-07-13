package cmds

import (
	"github.com/spf13/cobra"
	onessl "kubepack.dev/onessl/cmds"
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
