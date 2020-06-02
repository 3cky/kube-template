module github.com/3cky/kube-template

go 1.12

require (
	github.com/Masterminds/sprig/v3 v3.1.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/googleapis/gnostic v0.2.3-0.20180520015035-48a0ecefe2e4 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/onsi/ginkgo v1.12.3 // indirect
	github.com/pelletier/go-toml v1.2.1-0.20180724185102-c2dbbc24a979 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.5.1
	golang.org/x/net v0.0.0-20200528225125-3c3fba18258b // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/client-go v0.18.3
	k8s.io/kubernetes v1.18.3
	k8s.io/utils v0.0.0-20200529193333-24a76e807f40 // indirect
)

replace k8s.io/api => k8s.io/api v0.18.3

replace k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.18.3

replace k8s.io/apimachinery => k8s.io/apimachinery v0.18.4-rc.0

replace k8s.io/apiserver => k8s.io/apiserver v0.18.3

replace k8s.io/cli-runtime => k8s.io/cli-runtime v0.18.3

replace k8s.io/client-go => k8s.io/client-go v0.18.3

replace k8s.io/cloud-provider => k8s.io/cloud-provider v0.18.3

replace k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.18.3

replace k8s.io/code-generator => k8s.io/code-generator v0.18.4-rc.0

replace k8s.io/component-base => k8s.io/component-base v0.18.3

replace k8s.io/cri-api => k8s.io/cri-api v0.18.4-rc.0

replace k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.18.3

replace k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.18.3

replace k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.18.3

replace k8s.io/kube-proxy => k8s.io/kube-proxy v0.18.3

replace k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.18.3

replace k8s.io/kubectl => k8s.io/kubectl v0.18.3

replace k8s.io/kubelet => k8s.io/kubelet v0.18.3

replace k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.18.3

replace k8s.io/metrics => k8s.io/metrics v0.18.3

replace k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.18.3

replace k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.18.3

replace k8s.io/sample-controller => k8s.io/sample-controller v0.18.3
