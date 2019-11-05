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
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	"github.com/imdario/mergo"
	"github.com/spf13/cobra"
	"k8s.io/klog"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
)

func NewCmdMergeNodeConfig() *cobra.Command {
	joinCfg := &kubeadmapi.JoinConfiguration{}
	initCfg := &kubeadmapi.InitConfiguration{}

	fd := &kubeadmapi.FileDiscovery{}
	btd := &kubeadmapi.BootstrapTokenDiscovery{}

	var initCfgPath, joinCfgPath, token string
	cmd := &cobra.Command{
		Use:               "join-config",
		Short:             `Merge Kubeadm node configuration`,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			if len(fd.KubeConfigPath) != 0 {
				joinCfg.Discovery.File = fd
			} else {
				joinCfg.Discovery.BootstrapToken = btd
				if len(joinCfg.Discovery.BootstrapToken.Token) == 0 {
					joinCfg.Discovery.BootstrapToken.Token = token
				}
				if len(args) > 0 {
					if len(joinCfgPath) == 0 && len(args) > 1 {
						klog.Warningf("[join] WARNING: More than one API server endpoint supplied on command line %v. Using the first one.", args)
					}
					joinCfg.Discovery.BootstrapToken.APIServerEndpoint = args[0]
				}
			}
			if len(joinCfg.Discovery.TLSBootstrapToken) == 0 {
				joinCfg.Discovery.TLSBootstrapToken = token
			}

			if initCfgPath != "" {
				data, err := ioutil.ReadFile(initCfgPath)
				if err != nil {
					Fatal(err)
				}
				var in kubeadmapi.InitConfiguration
				err = yaml.Unmarshal(data, &in)
				if err != nil {
					Fatal(err)
				}

				err = mergo.MergeWithOverwrite(initCfg, in)
				if err != nil {
					Fatal(err)
				}
				if initCfg.LocalAPIEndpoint.AdvertiseAddress == "" {
					initCfg.LocalAPIEndpoint.AdvertiseAddress = "kube-apiserver"
				}
			}

			if joinCfgPath != "" {
				data, err := ioutil.ReadFile(joinCfgPath)
				if err != nil {
					Fatal(err)
				}
				var in kubeadmapi.JoinConfiguration
				err = yaml.Unmarshal(data, &in)
				if err != nil {
					Fatal(err)
				}

				err = mergo.Merge(joinCfg, in)
				if err != nil {
					Fatal(err)
				}
			}

			initData, err := yaml.Marshal(initCfg)
			if err != nil {
				Fatal(err)
			}
			fmt.Println(string(initData))

			joinCfg.APIVersion = "kubeadm.k8s.io/v1beta1"
			joinCfg.Kind = "JoinConfiguration"
			data, err := yaml.Marshal(joinCfg)
			if err != nil {
				Fatal(err)
			}
			fmt.Println("---")
			fmt.Println(string(data))
			os.Exit(0)
		},
	}
	// ref: https://github.com/kubernetes/kubernetes/blob/0b9efaeb34a2fc51ff8e4d34ad9bc6375459c4a4/cmd/kubeadm/app/cmd/join.go#L122
	cmd.PersistentFlags().StringVar(
		&joinCfgPath, "join-config", joinCfgPath,
		"Path to kubeadm config file")

	cmd.Flags().StringVar(&initCfgPath, "init-config", initCfgPath, "Path to kubeadm init config file (WARNING: Usage of a configuration file is experimental)")

	cmd.PersistentFlags().StringVar(
		&fd.KubeConfigPath, "discovery-file", "",
		"A file or url from which to load cluster information")

	cmd.PersistentFlags().StringVar(
		&btd.Token, "discovery-token", "",
		"A token used to validate cluster information fetched from the master")
	cmd.PersistentFlags().StringSliceVar(
		&btd.CACertHashes, "discovery-token-ca-cert-hash", []string{},
		"For token-based discovery, validate that the root CA public key matches this hash (format: \"<type>:<value>\").")
	cmd.PersistentFlags().BoolVar(
		&btd.UnsafeSkipCAVerification, "discovery-token-unsafe-skip-ca-verification", false,
		"For token-based discovery, allow joining without --discovery-token-ca-cert-hash pinning.")

	cmd.PersistentFlags().StringVar(
		&joinCfg.Discovery.TLSBootstrapToken, "tls-bootstrap-token", "",
		"A token used for TLS bootstrapping")

	cmd.PersistentFlags().StringVar(
		&joinCfg.NodeRegistration.Name, "node-name", joinCfg.NodeRegistration.Name,
		"Specify the node name.")
	cmd.PersistentFlags().StringVar(
		&joinCfg.NodeRegistration.CRISocket, "cri-socket", joinCfg.NodeRegistration.CRISocket,
		`Specify the CRI socket to connect to.`,
	)

	cmd.PersistentFlags().StringVar(
		&token, "token", "",
		"Use this token for both discovery-token and tls-bootstrap-token")

	return cmd
}
