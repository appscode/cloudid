package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/appscode/go/net"
	"github.com/appscode/mergo"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1alpha1"
	"k8s.io/kubernetes/cmd/kubeadm/app/features"
)

func NewCmdMergeMasterConfig() *cobra.Command {
	var (
		cfg = &kubeadmapi.MasterConfiguration{
			TokenTTL: &metav1.Duration{
				Duration: 0,
			},
		}
		sans []string
		isHa bool
	)
	var cfgPath string
	var etcdServerAddress string
	var featureGatesString string
	cmd := &cobra.Command{
		Use:               "master-config",
		Short:             `Merge Kubeadm master configuration`,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if cfg.FeatureGates, err = features.NewFeatureGate(&features.InitFeatureGates, featureGatesString); err != nil {
				os.Exit(1)
			}

			sanSet := sets.NewString(sans...)

			if cfgPath != "" {
				data, err := ioutil.ReadFile(cfgPath)
				if err != nil {
					Fatal(err)
				}
				var in kubeadmapi.MasterConfiguration
				err = yaml.Unmarshal(data, &in)
				if err != nil {
					Fatal(err)
				}
				sanSet.Insert(in.APIServerCertSANs...)

				err = mergo.Merge(cfg, in)
				if err != nil {
					Fatal(err)
				}
			}

			cfg.APIVersion = "kubeadm.k8s.io/v1alpha1"
			cfg.Kind = "MasterConfiguration"
			cfg.APIServerCertSANs = sanSet.List()
			if isHa {
				ips, _, err := net.RoutableIPs()
				if err != nil {
					Fatal(fmt.Errorf("failed to detect routable ips. Reason: %v", err))
				}
				if len(ips) == 0 {
					Fatal(fmt.Errorf("no routable ips found"))
				}
				nodeIp := ips[0]
				clusterType := "join"
				if etcdServerAddress == "" {
					clusterType = "seed"
				}
				extraArgs := map[string]string{
					"name":             cfg.NodeName,
					"cluster-type":     clusterType,
					"data-dir":         fmt.Sprintf("/var/lib/etcd/%v", cfg.NodeName),
					"listen-peer-urls": fmt.Sprintf("http://%s:2380", nodeIp),
					//"listen-metrics-urls":         "https://127.0.0.1:2381",
					"listen-client-urls":          "http://0.0.0.0:2379",
					"initial-advertise-peer-urls": fmt.Sprintf("http://%s:2380", nodeIp),
					"advertise-client-urls":       fmt.Sprintf("http://%s:2379", nodeIp),
					//"client-cert-auth":            "true",
					//"peer-client-cert-auth":       "false",
					"quota-backend-bytes": "2147483648",
					"v":              "10",
					"server-address": etcdServerAddress,
				}
				/*extraArgs := map[string]string{

				}*/
				cfg.Etcd.ExtraArgs = extraArgs
				//cfg.Etcd.ExtraArgs = append(cfg.Etcd.ExtraArgs, extraArgs...)
			}
			data, err := yaml.Marshal(cfg)
			if err != nil {
				Fatal(err)
			}
			fmt.Println(string(data))
			os.Exit(0)
		},
	}
	// ref: https://github.com/kubernetes/kubernetes/blob/0b9efaeb34a2fc51ff8e4d34ad9bc6375459c4a4/cmd/kubeadm/app/cmd/init.go#L141
	cmd.Flags().StringVar(
		&cfg.API.AdvertiseAddress, "apiserver-advertise-address", cfg.API.AdvertiseAddress,
		"The IP address the API Server will advertise it's listening on. 0.0.0.0 means the default network interface's address.",
	)
	cmd.Flags().Int32Var(
		&cfg.API.BindPort, "apiserver-bind-port", cfg.API.BindPort,
		"Port for the API Server to bind to",
	)
	cmd.Flags().StringVar(
		&cfg.Networking.ServiceSubnet, "service-cidr", cfg.Networking.ServiceSubnet,
		"Use alternative range of IP address for service VIPs",
	)
	cmd.Flags().StringVar(
		&cfg.Networking.PodSubnet, "pod-network-cidr", cfg.Networking.PodSubnet,
		"Specify range of IP addresses for the pod network; if set, the control plane will automatically allocate CIDRs for every node",
	)
	cmd.Flags().StringVar(
		&cfg.Networking.DNSDomain, "service-dns-domain", cfg.Networking.DNSDomain,
		`Use alternative domain for services, e.g. "myorg.internal"`,
	)
	cmd.Flags().StringVar(
		&cfg.KubernetesVersion, "kubernetes-version", cfg.KubernetesVersion,
		`Choose a specific Kubernetes version for the control plane`,
	)
	cmd.Flags().StringVar(
		&cfg.CertificatesDir, "cert-dir", cfg.CertificatesDir,
		`The path where to save and store the certificates`,
	)
	cmd.Flags().StringSliceVar(
		&sans, "apiserver-cert-extra-sans", sans,
		`Optional extra altnames to use for the API Server serving cert. Can be both IP addresses and dns names.`,
	)
	cmd.Flags().StringVar(
		&cfg.NodeName, "node-name", cfg.NodeName,
		`Specify the node name`,
	)
	cmd.Flags().StringVar(
		&cfg.Token, "token", cfg.Token,
		"The token to use for establishing bidirectional trust between nodes and masters.",
	)
	cmd.Flags().DurationVar(
		&cfg.TokenTTL.Duration, "token-ttl", cfg.TokenTTL.Duration,
		"The duration before the bootstrap token is automatically deleted. 0 means 'never expires'.",
	)
	cmd.Flags().StringVar(&featureGatesString, "feature-gates", featureGatesString, "A set of key=value pairs that describe feature gates for various features. "+
		"Options are:\n"+strings.Join(features.KnownFeatures(&features.InitFeatureGates), "\n"))
	cmd.Flags().StringVar(&cfgPath, "config", cfgPath, "Path to kubeadm config file (WARNING: Usage of a configuration file is experimental)")

	cmd.Flags().BoolVar(&isHa, "ha", false, "Enable to apply ha cluster")
	cmd.Flags().StringVar(&etcdServerAddress, "etcd-server", "", "Etcd server address to join member, example: http://127.0.0.1:2379")
	return cmd
}
