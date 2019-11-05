module pharmer.dev/pre-k

go 1.12

require (
	github.com/JamesClonk/vultr v2.0.1+incompatible
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/appscode/go v0.0.0-20191025021232-311ac347b3ef
	github.com/ghodss/yaml v1.0.0
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/imdario/mergo v0.3.5
	github.com/jpillora/go-ogle-analytics v0.0.0-20161213085824-14b04e0594ef
	github.com/konsorten/go-windows-terminal-sequences v1.0.2 // indirect
	github.com/pkg/errors v0.8.1
	github.com/sirupsen/logrus v1.4.1 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.3
	gomodules.xyz/cert v1.0.1
	k8s.io/apimachinery v0.0.0-20191025225532-af6325b3a843
	k8s.io/apiserver v0.0.0-20190516230822-f89599b3f645 // indirect
	k8s.io/cli-runtime v0.0.0-20190516231937-17bc0b7fcef5 // indirect
	k8s.io/client-go v11.0.0+incompatible
	k8s.io/cluster-bootstrap v0.0.0-20191025232351-410fafc3baf5 // indirect
	k8s.io/component-base v0.0.0-20190424053038-9fe063da3132 // indirect
	k8s.io/klog v1.0.0
	k8s.io/kubernetes v1.14.0
	kmodules.xyz/client-go v0.0.0-20191101042247-ee9566c9ac7f
	kubepack.dev/onessl v0.13.1
)

replace (
	cloud.google.com/go => cloud.google.com/go v0.34.0
	git.apache.org/thrift.git => github.com/apache/thrift v0.12.0
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v12.4.2+incompatible
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190315093550-53c4693659ed
	k8s.io/apimachinery => github.com/kmodules/apimachinery v0.0.0-20190508045248-a52a97a7a2bf
	k8s.io/apiserver => github.com/kmodules/apiserver v0.0.0-20190811223248-5a95b2df4348
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190314001948-2899ed30580f
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190314002645-c892ea32361a
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190314000054-4a91899592f4
	k8s.io/klog => k8s.io/klog v0.3.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190314000639-da8327669ac5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30
	k8s.io/metrics => k8s.io/metrics v0.0.0-20190314001731-1bd6a4002213
	k8s.io/utils => k8s.io/utils v0.0.0-20190514214443-0a167cbac756
	sigs.k8s.io/structured-merge-diff => sigs.k8s.io/structured-merge-diff v0.0.0-20190302045857-e85c7b244fd2
)
