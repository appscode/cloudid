## pre-k merge config

Merge Kubeadm initial configuration

### Synopsis

Merge Kubeadm initial configuration

```
pre-k merge config [flags]
```

### Options

```
      --apiserver-advertise-address string   The IP address the API Server will advertise it's listening on. 0.0.0.0 means the default network interface's address.
      --apiserver-bind-port int32            Port for the API Server to bind to (default 6443)
      --apiserver-cert-extra-sans strings    Optional extra altnames to use for the API Server serving cert. Can be both IP addresses and dns names.
      --cert-dir string                      The path where to save and store the certificates (default "/etc/kubernetes/pki")
      --cluster-config string                Path to kubeadm cluster config file (WARNING: Usage of a configuration file is experimental)
      --cri-socket string                    Specify the CRI socket to connect to.
      --etcd-server string                   Etcd server address to join member, example: 127.0.0.1
      --feature-gates string                 A set of key=value pairs that describe feature gates for various features. Options are:
                                             
      --ha                                   Enable to apply ha cluster
  -h, --help                                 help for config
      --init-config string                   Path to kubeadm init config file (WARNING: Usage of a configuration file is experimental)
      --kubernetes-version string            Choose a specific Kubernetes version for the control plane (default "stable-1")
      --node-name string                     Specify the node name
      --pod-network-cidr string              Specify range of IP addresses for the pod network; if set, the control plane will automatically allocate CIDRs for every node
      --service-cidr string                  Use alternative range of IP address for service VIPs (default "10.96.0.0/12")
      --service-dns-domain string            Use alternative domain for services, e.g. "myorg.internal" (default "cluster.local")
      --tls-enabled                          Enable tls to secure etcd (default true)
      --token string                         The token to use for establishing bidirectional trust between nodes and masters.
      --token-ttl duration                   The duration before the bootstrap token is automatically deleted. 0 means 'never expires'.
```

### Options inherited from parent commands

```
      --alsologtostderr                  log to standard error as well as files
      --analytics                        Send analytical events to Google Guard (default true)
      --log-flush-frequency duration     Maximum number of seconds between log flushes (default 5s)
      --log_backtrace_at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log_dir string                   If non-empty, write log files in this directory
      --logtostderr                      log to standard error instead of files
      --stderrthreshold severity         logs at or above this threshold go to stderr
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### SEE ALSO

* [pre-k merge](pre-k_merge.md)	 - Merge Kubeadm config

