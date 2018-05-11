package cmds

import (
	"io/ioutil"

	"github.com/appscode/go/term"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm"
	kubeadmphase "k8s.io/kubernetes/cmd/kubeadm/app/phases/certs"
)

func NewCmdEtcd() *cobra.Command {
	var cfgPath string
	cmd := &cobra.Command{
		Use:               "etcd",
		Short:             "Configure etcd",
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
			if conf.CertificatesDir != "" {
				if err := generateCerts(&conf); err != nil {
					term.Fatalln(err)
				}
			}

		},
	}
	cmd.Flags().StringVar(&cfgPath, "config", cfgPath, "Path to kubeadm config file (WARNING: Usage of a configuration file is experimental)")

	return cmd
}

func generateCerts(cfg *kubeadmapi.MasterConfiguration) error {
	if err := kubeadmphase.CreateEtcdServerCertAndKeyFiles(cfg); err != nil {
		return err
	}

	if err := kubeadmphase.CreateEtcdPeerCertAndKeyFiles(cfg); err != nil {
		return err
	}

	if err := kubeadmphase.CreateEtcdHealthcheckClientCertAndKeyFiles(cfg); err != nil {
		return err
	}
	if err := kubeadmphase.CreateAPIServerEtcdClientCertAndKeyFiles(cfg); err != nil {
		return err
	}
	return nil
}
