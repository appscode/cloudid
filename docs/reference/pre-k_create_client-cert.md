## pre-k create client-cert

Generate client certificate pair

### Synopsis

Generate client certificate pair

```
pre-k create client-cert [flags]
```

### Options

```
      --cert-dir string        Path to directory where pki files are stored. (default "/etc/kubernetes/pki")
  -h, --help                   help for client-cert
  -o, --organization strings   Name of client organizations.
      --overwrite              Overwrite existing cert/key pair
  -p, --prefix string          Prefix added to certificate files
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

* [pre-k create](pre-k_create.md)	 - create PKI

