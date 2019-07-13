package cmds

import (
	"github.com/spf13/cobra"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
	onessl "kubepack.dev/onessl/cmds"
)

func NewCmdCreate() *cobra.Command {
	var (
		certDir = kubeadmapi.DefaultCertificatesDir
	)
	cmd := &cobra.Command{
		Use:               "create",
		Short:             `create PKI`,
		DisableAutoGenTag: true,
	}
	cmd.AddCommand(onessl.NewCmdCreateCA(certDir))
	cmd.AddCommand(onessl.NewCmdCreateServer(certDir))
	cmd.AddCommand(onessl.NewCmdCreatePeer(certDir))
	cmd.AddCommand(onessl.NewCmdCreateClient(certDir))
	return cmd
}
