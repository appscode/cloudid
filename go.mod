module pharmer.dev/pre-k

go 1.12

require (
	github.com/JamesClonk/vultr v0.0.0-20180101102256-fa1c0367800d
	github.com/appscode/go v0.0.0-20190621064509-6b292c9166e3
	github.com/ghodss/yaml v1.0.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/imdario/mergo v0.3.7
	github.com/jpillora/go-ogle-analytics v0.0.0-20161213085824-14b04e0594ef
	github.com/juju/ratelimit v0.0.0-20151125201925-77ed1c8a0121 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	gomodules.xyz/cert v1.0.0
	k8s.io/apimachinery v0.0.0-20190424052434-11f1676e3da4
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/cluster-bootstrap v0.0.0-20190424054052-c2758412356a // indirect
	k8s.io/klog v0.3.0
	k8s.io/kubernetes v1.14.0
	k8s.io/utils v0.0.0-20180726175726-66066c83e385 // indirect
	kmodules.xyz/client-go v0.0.0-20190715080709-7162a6c90b04
	kubepack.dev/onessl v0.13.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest/autorest v0.5.0
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190315093550-53c4693659ed
	k8s.io/apimachinery => github.com/kmodules/apimachinery v0.0.0-20190508045248-a52a97a7a2bf
	k8s.io/apiserver => github.com/kmodules/apiserver v0.0.0-20190508082252-8397d761d4b5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190314001948-2899ed30580f
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190314002645-c892ea32361a
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190314000054-4a91899592f4
	k8s.io/klog => k8s.io/klog v0.3.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190314000639-da8327669ac5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30
	k8s.io/metrics => k8s.io/metrics v0.0.0-20190314001731-1bd6a4002213
	k8s.io/utils => k8s.io/utils v0.0.0-20190221042446-c2654d5206da
)
