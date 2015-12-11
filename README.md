kube-template
===============

Inspired by HashiCorp's [Consul Template][], `kube-template` is a utility that queries a [Kubernetes][] API server and uses its objects as input data for specified set of templates. Also it can run arbitrary commands for templates with updated output.

**Please note this is work in progress.**

Usage
-----

### Options

```
      --alsologtostderr[=false]: log to standard error as well as files
  -c, --config="": config file (default is ./kube-template.(yaml|json))
      --dry-run[=false]: don't write template output, dump result to stdout
      --help-md[=false]: get help in Markdown format
      --log-backtrace-at=:0: when logging hits line file:N, emit a stack trace
      --log-dir="": If non-empty, write log files in this directory
      --log-flush-frequency=5s: Maximum number of seconds between log flushes
      --logtostderr[=true]: log to standard error instead of files
      --once[=false]: run template processing once and exit
  -p, --poll-time=15s: Kubernetes API server poll time
  -s, --server="": the address and port of the Kubernetes API server
      --stderrthreshold=2: logs at or above this threshold go to stderr
  -t, --template=[]: adds a new template to watch on disk in the format
		'templatePath:outputPath[:command]'. This option is additive
		and may be specified multiple times for multiple templates.
      --v=0: log level for V logs
      --vmodule=: comma-separated list of pattern=N settings for file-filtered logging
```

### Command Line

Process single template using data from remote Kubernetes API server and exit:

```shell
$ kube-template \
    --server http://kube-api.server.com:8080 \
    --template "/tmp/input.tmpl:/tmp/output.txt" \ 
    --once 
```

Monitor local Kubernetes API server for updates every 30 seconds and update nginx and haproxy configuration files with reload:

```shell
$ kube-template \
    --template "/tmp/nginx.tmpl:/etc/nginx/nginx.conf:service nginx reload" \ 
    --template "/tmp/haproxy.tmpl:/etc/haproxy/haproxy.conf:service haproxy reload" \
    --poll-time=30s
```

### Configuration File

`kube-template` looks for `kube-template.json` or `kube-template.yaml` configuration file in current working directory or file name specified by `--config` command line option.
 
 YAML configuration file example:
 
```yaml
 server: http://localhost:8080
 poll-time: 10s
 
 templates:
   - path: in1.tmpl
     output: out1.txt
     command: action1.sh
   - path: in1.tmpl
     output: out2.txt
   - path: incomplete1.tmpl
```

___Please note___: templates specified on the command line take precedence over those defined in a config file.

### Templating Language

`kube-template` works with templates in the [Go Template][] format. In addition to the [standard template functions][Go Template], `kube-template` provides the following functions:

#### Kubernetes API

##### `pods`
```
{{pods "selector" "namespace"}}
```
Query Kubernetes API server for [pods](https://github.com/kubernetes/kubernetes/blob/master/docs/user-guide/pods.md) from given `namespace` (`default` if not specified) matching given `selector` (empty to get all pods).
 
Example:
```
{{range pods}}
{{.Name}}: {{.Status.Phase}}
{{end}}
```

##### `services`
```
{{services "selector" "namespace"}}
```
Query Kubernetes API server for [services](https://github.com/kubernetes/kubernetes/blob/master/docs/user-guide/services.md) from given `namespace` (`default` if not specified) matching given `selector` (empty to get all services). 

##### `replicationcontrollers`
```
{{replicationcontrollers "selector" "namespace"}}
```
Query Kubernetes API server for [replication controllers](https://github.com/kubernetes/kubernetes/blob/master/docs/user-guide/replication-controller.md) from given `namespace` (`default` if not specified) matching given `selector` (empty to get all replication controllers). 

##### `events`
```
{{events "selector" "namespace"}}
```
Query Kubernetes API server for events from given `namespace` (`default` if not specified) matching given `selector` (empty to get all events). 

##### `endpoints`
```
{{endpoints "selector" "namespace"}}
```
Query Kubernetes API server for endpoints from given `namespace` (`default` if not specified) matching given `selector` (empty to get all endpoints). 

##### `nodes`
```
{{nodes "selector"}}
```
Query Kubernetes API server for nodes matching given `selector` (empty to get all nodes). 

##### `namespaces`
```
{{namespaces "selector"}}
```
Query Kubernetes API server for namespaces matching given `selector` (empty to get all namespaces). 

- - -

#### Helper Functions

##### `add`
```
{{add a b}}
```
Returns the sum of two integers, `a` and `b`.

##### `sub`
```
{{sub a b}}
```
Returns the subtract of integer `b` from integer `a`.

### EXAMPLES
**To be done.**

### TODO
* Track Kubernetes changes using resource watch API
* More API and helper template functions
* Add options for API server authentication
* Implement configuration file reloading
* Add tests

## Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## License

`kube-template` is released under the Apache 2.0 license. See [LICENSE.txt](https://github.com/3cky/kube-template/blob/master/LICENSE.txt)

[Kubernetes]: http://kubernetes.io/ "Manage a cluster of Linux containers as a single system to accelerate Dev and simplify Ops"
[Consul Template]: https://github.com/hashicorp/consul-template "A convenient way to populate values from Consul into the filesystem using the consul-template daemon"
[Go Template]: http://golang.org/pkg/text/template/ "Go Template"
