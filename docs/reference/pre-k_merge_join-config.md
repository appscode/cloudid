## pre-k merge join-config

Merge Kubeadm node configuration

### Synopsis

Merge Kubeadm node configuration

```
pre-k merge join-config [flags]
```

### Options

```
      --cri-socket string                             Specify the CRI socket to connect to.
      --discovery-file string                         A file or url from which to load cluster information
      --discovery-token string                        A token used to validate cluster information fetched from the master
      --discovery-token-ca-cert-hash strings          For token-based discovery, validate that the root CA public key matches this hash (format: "<type>:<value>").
      --discovery-token-unsafe-skip-ca-verification   For token-based discovery, allow joining without --discovery-token-ca-cert-hash pinning.
  -h, --help                                          help for join-config
      --init-config string                            Path to kubeadm init config file (WARNING: Usage of a configuration file is experimental)
      --join-config string                            Path to kubeadm config file
      --node-name string                              Specify the node name.
      --tls-bootstrap-token string                    A token used for TLS bootstrapping
      --token string                                  Use this token for both discovery-token and tls-bootstrap-token
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

