module github.com/open-cluster-management/multicloud-operators-application

require (
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/kubernetes-sigs/application v0.8.1
	github.com/onsi/ginkgo v1.12.0
	github.com/onsi/gomega v1.9.0
	github.com/open-cluster-management/multicloud-operators-deployable v0.0.0-20200603180154-d1d17d718c30
	github.com/open-cluster-management/multicloud-operators-subscription v0.0.0-20200608164921-3102dd215933
	github.com/operator-framework/operator-sdk v0.18.1
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.9.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.6.0
	k8s.io/api v0.18.2
	k8s.io/apiextensions-apiserver v0.18.2
	k8s.io/apimachinery v0.18.2
	k8s.io/client-go v13.0.0+incompatible
	k8s.io/klog v1.0.0
	sigs.k8s.io/controller-runtime v0.6.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v13.3.2+incompatible // Required by OLM
	golang.org/x/text => golang.org/x/text v0.3.5
	k8s.io/client-go => k8s.io/client-go v0.18.2 // Required by prometheus-operator
)

go 1.15
