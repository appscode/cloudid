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
	"os"
	"time"

	"github.com/appscode/go/log"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	core_util "kmodules.xyz/client-go/core/v1"
)

func NewCmdCheckMasterStatus() *cobra.Command {
	var (
		masterURL      string
		kubeconfigPath string
		interval       = 2 * time.Second
		timeout        time.Duration
	)
	cmd := &cobra.Command{
		Use:               "master-status",
		Short:             "Checks whether master(s) are running and ready",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
			if err != nil {
				log.Fatalln(err)
			}
			kubeClient := kubernetes.NewForConfigOrDie(config)

			check := func() (done bool, err error) {
				nodes, err := kubeClient.CoreV1().Nodes().List(metav1.ListOptions{
					LabelSelector: labels.SelectorFromSet(map[string]string{
						"node-role.kubernetes.io/master": "",
					}).String(),
				})
				if err != nil {
					return false, nil
				}
				for _, node := range nodes.Items {
					if !core_util.NodeReady(node) {
						return false, nil
					}
				}
				return len(nodes.Items) > 0, nil // there must be at least one master
			}

			if timeout == 0 {
				ok, _ := check()
				if !ok {
					os.Exit(1)
				}
			} else if timeout < 0 {
				err = wait.PollImmediateInfinite(interval, check)
				if err != nil {
					os.Exit(1)
				}
			} else {
				err = wait.PollImmediate(interval, timeout, check)
				if err != nil {
					os.Exit(1)
				}
			}
		},
	}
	cmd.Flags().StringVar(&masterURL, "master", masterURL, "The address of the Kubernetes API server (overrides any value in kubeconfig)")
	cmd.Flags().StringVar(&kubeconfigPath, "kubeconfig", kubeconfigPath, "Path to kubeconfig file with authorization information (the master location is set by the master flag).")
	cmd.Flags().DurationVar(&interval, "interval", interval, "Interval between checks")
	cmd.Flags().DurationVar(&timeout, "timeout", timeout, "Timeout for check master status")

	return cmd
}
