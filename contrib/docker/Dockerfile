FROM alpine:latest

MAINTAINER: me@foo.org

COPY bin/kube-template /usr/bin/kube-template
COPY conf/kube-template.yaml /usr/bin/kube-template.yaml
COPY conf/in.txt.tmpl /etc/kube-template/in.txt.tmpl

RUN set -o pipefail && \
addgroup -g 987 kube  && \
adduser -h /home/kube -u 990 -G kube -g "Kubernetes user" -S -D kube  && \
mkdir -m 775 /var/lib/kube-template  && \
chown kube: /var/lib/kube-template  && \
chown -R kube: /etc/kube-template  && \
chmod 775 /etc/kube-template  && \
chmod 640 /etc/kube-template/in.txt.tmpl  && \
chown kube: /usr/bin/kube-template.yaml  && \
chmod 640 /usr/bin/kube-template.yaml  && \
chmod 755 /usr/bin/kube-template

WORKDIR /usr/bin
USER kube

CMD ["--guess-kube-api-settings"]
ENTRYPOINT ["/usr/bin/kube-template"]
