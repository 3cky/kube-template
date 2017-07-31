kube-template
===============

Inspired by HashiCorp's [Consul Template][], `kube-template` is a utility that queries a [Kubernetes][] API server and uses its objects as input data for specified set of templates. Also it can run arbitrary commands for templates with updated output.

**Please note this is work in progress.**

Installation
------------

- Install [Go vendor tool][]: `go get -u github.com/kardianos/govendor`
- For Go 1.5 make sure you have set **GO15VENDOREXPERIMENT=1** environment variable (set by default in Go 1.6)
- `git pull https://github.com/3cky/kube-template.git`
- `cd kube-template`
- `make install`

Docker

If you want to create a Docker image
- adjust Dockerfile in contrib/docker directory to fit your needs
- configure settings files in contrib/docker/conf directory
- issue `make docker`

Note - the built image is already configured to automatically locate Kubernetes API along with everything is needed to connect (e.g. Certificate Authority file, security token and so on). 

Usage
-----

### Options

```
      --alsologtostderr                  log to standard error as well as files
  -c, --config string                    config file (default is ./kube-template.(yaml|json))
      --dry-run                          don't write template output, dump result to stdout
      --guess-kube-api-settings          guess Kubernetes API settings from POD environment
      --help-md                          get help in Markdown format
      --log-backtrace-at traceLocation   when logging hits line file:N, emit a stack trace (default :0)
      --log-dir string                   If non-empty, write log files in this directory
      --log-flush-frequency duration     Maximum number of seconds between log flushes (default 5s)
      --logtostderr                      log to standard error instead of files (default true)
      --master string                    Kubernetes API server address (default is http://127.0.0.1:8080/)
      --once                             run template processing once and exit
  -p, --poll-time duration               Kubernetes API server poll time (0 disables server polling) (default 15s)
      --stderrthreshold severity         logs at or above this threshold go to stderr (default 2)
  -t, --template stringSlice             adds a new template to watch on disk in the format
		'templatePath:outputPath[:command]'. This option is additive
		and may be specified multiple times for multiple templates
  -v, --v Level                          log level for V logs
      --vmodule moduleSpec               comma-separated list of pattern=N settings for file-filtered logging
```

### Command Line

Process single template using data from remote Kubernetes API server and exit:

```shell
$ kube-template \
    --master="http://kube-api.server.com:8080" \
    --template="/tmp/input.tmpl:/tmp/output.txt" \ 
    --once 
```

Monitor local Kubernetes API server for updates every 30 seconds and update nginx and haproxy configuration files with reload:

```shell
$ kube-template \
    --template="/tmp/nginx.tmpl:/etc/nginx/nginx.conf:service nginx reload" \ 
    --template="/tmp/haproxy.tmpl:/etc/haproxy/haproxy.conf:service haproxy reload" \
    --poll-time=30s
```

### Configuration File

`kube-template` looks for `kube-template.json` or `kube-template.yaml` configuration file in current working directory or file name specified by `--config` command line option.
 
 YAML configuration file example:
 
```yaml
 master: http://localhost:8080
 poll-time: 10s
 
 templates:
   - path: in.txt.tmpl
     output: out.txt
     command: action.sh

   - path: in.html.tmpl
     output: out.html
```

___Please note___: templates specified on the command line take precedence over those defined in a config file.

### Signals

- **TERM, QUIT, INT:** graceful shutdown
- **HUP:** reload configuration file

### Templating Language

`kube-template` works with templates in the [Go Template][] format. In addition to the [standard template functions][Go Template], `kube-template` provides the following functions:

#### Kubernetes API

##### `pods`
```
{{pods "selector" "namespace"}}
```
Query Kubernetes API server for [pods](https://kubernetes.io/docs/user-guide/pods/) from given `namespace` (`default` if not specified) matching given `selector` (empty to get all pods).
 
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
Query Kubernetes API server for [services](https://kubernetes.io/docs/user-guide/services/) from given `namespace` (`default` if not specified) matching given `selector` (empty to get all services).

##### `replicationcontrollers`
```
{{replicationcontrollers "selector" "namespace"}}
```
Query Kubernetes API server for [replication controllers](https://kubernetes.io/docs/user-guide/replication-controller/) from given `namespace` (`default` if not specified) matching given `selector` (empty to get all replication controllers).

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

##### `mul`
```
{{mul a b}}
```
Returns the multiplication of integer `a` by integer `b`.

### EXAMPLES
Populate nginx upstream with pods labeled 'name=test-pod':
```
upstream test-pod-upstream {
{{with $pods := pods "name=test-pod"}}
    ip_hash;
    {{range $pods}}
    server {{.Status.PodIP}}:8080;
    {{end}}
{{else}}
    # No pods with given label were found, use stub server
    server 127.0.0.1:8081;
{{end}}
}
```

Create haproxy TCP balancer for pods listening on container port 9000:
```
{{with $pods := pods "name=test-pod"}}
{{$maxconn := 1000}}
listen test-pod-balancer
    bind *:9000
    mode tcp
    maxconn {{mul $maxconn (len $pods)}}
    balance roundrobin
    {{range $pods}}
    server {{.Name}} {{.Status.PodIP}}:9000 maxconn {{$maxconn}} check inter 5000 rise 3 fall 3
    {{end}}
    server stub 127.0.0.1:7690 backup
{{end}}
```

### TODO
* Track Kubernetes changes using resource watch API
* More API and helper template functions
* Add options for API server authentication
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
[Go vendor tool]: https://github.com/kardianos/govendor "Go vendor tool that works with the standard vendor file."
