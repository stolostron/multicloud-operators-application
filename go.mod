module github.com/open-cluster-management/multicloud-operators-application

go 1.16

require (
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/kubernetes-sigs/application v0.8.1
	github.com/onsi/ginkgo v1.16.4
	github.com/onsi/gomega v1.13.0
	github.com/open-cluster-management/multicloud-operators-deployable v1.2.4-0-20210816-f9fe854
	github.com/open-cluster-management/multicloud-operators-subscription v1.2.4-0-20210817-7443bc9
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.0
	k8s.io/api v0.21.3
	k8s.io/apiextensions-apiserver v0.21.3
	k8s.io/apimachinery v0.21.3
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.9.1
)

replace (
	github.com/ulikunitz/xz => github.com/ulikunitz/xz v0.5.10
	k8s.io/api => k8s.io/api v0.19.3
	k8s.io/client-go => k8s.io/client-go v0.19.3
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.6.3
)
