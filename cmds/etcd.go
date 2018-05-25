package cmds

import (
	"io/ioutil"

	"github.com/appscode/go/term"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmconstants "k8s.io/kubernetes/cmd/kubeadm/app/constants"
	kubeadmphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
)

func NewCmdCreateEtcdCerts() *cobra.Command {
	var cfgPath string
	cmd := &cobra.Command{
		Use:               "etcd-certs",
		Short:             "Create etcd client and peer certificate pairs",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if cfgPath == "" {
				term.Fatalln("Master configuration file must not be empty")
			}
			data, err := ioutil.ReadFile(cfgPath)
			if err != nil {
				Fatal(err)
			}
			var conf kubeadmapi.MasterConfiguration
			err = yaml.Unmarshal(data, &conf)
			if err != nil {
				Fatal(err)
			}
			if conf.CertificatesDir == "" {
				conf.CertificatesDir = kubeadmconstants.KubernetesDir + "/pki"
			}

			if err := generateCerts(&conf); err != nil {
				term.Fatalln(err)
			}

		},
	}
	cmd.Flags().StringVar(&cfgPath, "config", cfgPath, "Path to kubeadm config file (WARNING: Usage of a configuration file is experimental)")

	return cmd
}

func generateCerts(cfg *kubeadmapi.MasterConfiguration) error {
	certActions := []func(cfg *kubeadmapi.MasterConfiguration) error{
		kubeadmphase.CreateEtcdServerCertAndKeyFiles,
		kubeadmphase.CreateEtcdPeerCertAndKeyFiles,
		kubeadmphase.CreateEtcdHealthcheckClientCertAndKeyFiles,
		kubeadmphase.CreateAPIServerEtcdClientCertAndKeyFiles,
	}

	for _, action := range certActions {
		err := action(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}
