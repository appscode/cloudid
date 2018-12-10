package cmds

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/appscode/mergo"
	"github.com/ghodss/yaml"
	"github.com/spf13/cobra"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/sets"
	kubeadmapi "k8s.io/kubernetes/cmd/kubeadm/app/apis/kubeadm/v1beta1"
	"k8s.io/kubernetes/cmd/kubeadm/app/features"
	kubeadmutil "k8s.io/kubernetes/cmd/kubeadm/app/util"
)

func NewCmdMergeMasterConfig() *cobra.Command {
	var (
		initCfg    = &kubeadmapi.InitConfiguration{}
		clusterCfg = &kubeadmapi.ClusterConfiguration{}
		sans       []string
		isHa       bool
		tlsEnabled bool
		token      string
		tokenTTL   time.Duration = 0
	)
	var initCfgPath string
	var clusterCfgPath string

	var etcdServerAddress string
	var featureGatesString string

	kubeadmapi.SetDefaults_InitConfiguration(initCfg)
	kubeadmapi.SetDefaults_ClusterConfiguration(clusterCfg)

	cmd := &cobra.Command{
		Use:               "config",
		Short:             `Merge Kubeadm initial configuration`,
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			var err error
			if clusterCfg.FeatureGates, err = features.NewFeatureGate(&features.InitFeatureGates, featureGatesString); err != nil {
				os.Exit(1)
			}
			if token != "" {
				bt, err := kubeadmapi.NewBootstrapTokenString(token)
				if err != nil {
					Fatal(err)
				}
				initCfg.BootstrapTokens = []kubeadmapi.BootstrapToken{
					{
						Token: bt,
						TTL: &metav1.Duration{
							tokenTTL,
						},
					},
				}
			}

			sanSet := sets.NewString(sans...)

			kubeadmutil.CheckErr(err)

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
			}
			if clusterCfgPath != "" {
				data, err := ioutil.ReadFile(clusterCfgPath)
				if err != nil {
					Fatal(err)
				}
				var in kubeadmapi.ClusterConfiguration
				err = yaml.Unmarshal(data, &in)
				if err != nil {
					Fatal(err)
				}

				sanSet.Insert(in.APIServer.CertSANs...)
				err = mergo.MergeWithOverwrite(clusterCfg, in)
				if err != nil {
					Fatal(err)
				}
			}

			clusterCfg.APIServer.CertSANs = sanSet.List()

			initData, err := yaml.Marshal(initCfg)
			if err != nil {
				Fatal(err)
			}
			fmt.Println(string(initData))

			clusterData, err := yaml.Marshal(clusterCfg)
			if err != nil {
				Fatal(err)
			}
			fmt.Println("---")
			fmt.Println(string(clusterData))
			os.Exit(0)
		},
	}
	// ref: https://github.com/kubernetes/kubernetes/blob/0b9efaeb34a2fc51ff8e4d34ad9bc6375459c4a4/cmd/kubeadm/app/cmd/init.go#L141
	cmd.Flags().StringVar(
		&initCfg.LocalAPIEndpoint.AdvertiseAddress, "apiserver-advertise-address", initCfg.LocalAPIEndpoint.AdvertiseAddress,
		"The IP address the API Server will advertise it's listening on. 0.0.0.0 means the default network interface's address.",
	)
	cmd.Flags().Int32Var(
		&initCfg.LocalAPIEndpoint.BindPort, "apiserver-bind-port", initCfg.LocalAPIEndpoint.BindPort,
		"Port for the API Server to bind to",
	)
	cmd.Flags().StringVar(
		&initCfg.Networking.ServiceSubnet, "service-cidr", initCfg.Networking.ServiceSubnet,
		"Use alternative range of IP address for service VIPs",
	)
	cmd.Flags().StringVar(
		&initCfg.Networking.PodSubnet, "pod-network-cidr", initCfg.Networking.PodSubnet,
		"Specify range of IP addresses for the pod network; if set, the control plane will automatically allocate CIDRs for every node",
	)
	cmd.Flags().StringVar(
		&initCfg.Networking.DNSDomain, "service-dns-domain", initCfg.Networking.DNSDomain,
		`Use alternative domain for services, e.g. "myorg.internal"`,
	)
	cmd.Flags().StringVar(
		&clusterCfg.KubernetesVersion, "kubernetes-version", clusterCfg.KubernetesVersion,
		`Choose a specific Kubernetes version for the control plane`,
	)
	cmd.Flags().StringVar(
		&clusterCfg.CertificatesDir, "cert-dir", clusterCfg.CertificatesDir,
		`The path where to save and store the certificates`,
	)
	cmd.Flags().StringSliceVar(
		&sans, "apiserver-cert-extra-sans", sans,
		`Optional extra altnames to use for the API Server serving cert. Can be both IP addresses and dns names.`,
	)
	cmd.Flags().StringVar(
		&initCfg.NodeRegistration.Name, "node-name", initCfg.NodeRegistration.Name,
		`Specify the node name`,
	)
	cmd.Flags().StringVar(
		&initCfg.NodeRegistration.CRISocket, "cri-socket", initCfg.NodeRegistration.CRISocket,
		`Specify the CRI socket to connect to.`,
	)
	cmd.Flags().StringVar(
		&token, "token", token,
		"The token to use for establishing bidirectional trust between nodes and masters.",
	)
	cmd.Flags().DurationVar(
		&tokenTTL, "token-ttl", tokenTTL,
		"The duration before the bootstrap token is automatically deleted. 0 means 'never expires'.",
	)
	cmd.Flags().StringVar(&featureGatesString, "feature-gates", featureGatesString, "A set of key=value pairs that describe feature gates for various features. "+
		"Options are:\n"+strings.Join(features.KnownFeatures(&features.InitFeatureGates), "\n"))
	cmd.Flags().StringVar(&initCfgPath, "init-config", initCfgPath, "Path to kubeadm init config file (WARNING: Usage of a configuration file is experimental)")
	cmd.Flags().StringVar(&clusterCfgPath, "cluster-config", clusterCfgPath, "Path to kubeadm cluster config file (WARNING: Usage of a configuration file is experimental)")

	cmd.Flags().BoolVar(&isHa, "ha", false, "Enable to apply ha cluster")
	cmd.Flags().StringVar(&etcdServerAddress, "etcd-server", "", "Etcd server address to join member, example: 127.0.0.1")
	cmd.Flags().BoolVar(&tlsEnabled, "tls-enabled", true, "Enable tls to secure etcd")
	return cmd
}
