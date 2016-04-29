Assumptions of the service file
-------------------------------

1. kube-template binary live in /usr/bin.
2. kube-template is run as root.
3. Configuration is done in via environment file in `/etc/kube-template/`.

Defaults in the environment file
--------------------------------
1. Default to log to stdout/journald instead of directly to disk, see: [KUBE_LOGTOSTDERR](config)
2. Kubernetes API server to connect is on http://127.0.0.1:8080, see: [KUBE_MASTER](config)
3. Template configurations are defined in `/etc/kube-template/kube-template.yaml`, see: [KUBE_TEMPLATE_ARGS](config)

How to use these files
----------------------

Place service file to `/etc/systemd/system/` directory, environment file (`config`) and your custom `kube-template.yaml` to `/etc/kube-template/` directory.

Start kube-template service: `systemctl start kube-template`.

Reload kube-template config: `systemctl reload kube-template`.

Enable kube-template service to run on boot: `systemctl enable kube-template`.
