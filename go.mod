module github.com/stolostron/multicloud-operators-application

go 1.19

require (
	github.com/ghodss/yaml v1.0.1-0.20190212211648-25d852aebe32
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.20.0
	github.com/open-cluster-management/multicloud-operators-deployable v1.2.4-1-20220201-2d1add0
	github.com/open-cluster-management/multicloud-operators-subscription v1.2.4-0-20211122-7277a37
	github.com/pkg/errors v0.9.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.8.0
	k8s.io/api v0.24.3
	k8s.io/apiextensions-apiserver v0.24.3
	k8s.io/apimachinery v0.24.3
	k8s.io/client-go v12.0.0+incompatible
	k8s.io/klog v1.0.0
	sigs.k8s.io/application v0.8.3
	sigs.k8s.io/controller-runtime v0.11.2
)

require (
	github.com/PuerkitoBio/purell v1.1.1 // indirect
	github.com/PuerkitoBio/urlesc v0.0.0-20170810143723-de5bf2ad4578 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful v2.16.0+incompatible // indirect
	github.com/evanphx/json-patch v4.12.0+incompatible // indirect
	github.com/fsnotify/fsnotify v1.5.1 // indirect
	github.com/go-logr/logr v1.2.0 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.19.5 // indirect
	github.com/go-openapi/spec v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.14 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/google/gofuzz v1.1.0 // indirect
	github.com/google/uuid v1.1.2 // indirect
	github.com/googleapis/gnostic v0.5.5 // indirect
	github.com/imdario/mergo v0.3.12 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.6 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/open-cluster-management/api v0.0.0-20210513122330-d76f10481f05 // indirect
	github.com/open-cluster-management/multicloud-operators-channel v1.2.5-0-20211122-79abb30 // indirect
	github.com/open-cluster-management/multicloud-operators-placementrule v1.2.4-0-20211122-be034 // indirect
	github.com/open-cluster-management/multicloud-operators-subscription-release v1.2.4-0-20211122-8309641 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_golang v1.12.1 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.32.1 // indirect
	github.com/prometheus/procfs v0.7.3 // indirect
	golang.org/x/net v0.0.0-20220725212005-46097bf591d3 // indirect
	golang.org/x/oauth2 v0.0.0-20210819190943-2bc19b11175f // indirect
	golang.org/x/sys v0.0.0-20220722155257-8c9f86f7a55f // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	golang.org/x/text v0.3.8 // indirect
	golang.org/x/time v0.0.0-20210723032227-1f47c861a9ac // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/tomb.v1 v1.0.0-20141024135613-dd632973f1e7 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/component-base v0.24.3 // indirect
	k8s.io/klog/v2 v2.60.1 // indirect
	k8s.io/kube-openapi v0.0.0-20220328201542-3ee0da9b0b42 // indirect
	k8s.io/utils v0.0.0-20220210201930-3a6ce19ff2f9 // indirect
	sigs.k8s.io/json v0.0.0-20211208200746-9f7c6b3444d2 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
	sigs.k8s.io/yaml v1.3.0 // indirect
)

replace (
	github.com/open-cluster-management/api => open-cluster-management.io/api v0.0.0-20210513122330-d76f10481f05
	github.com/open-cluster-management/multicloud-operators-channel => github.com/stolostron/multicloud-operators-channel v1.2.5-0-20211122-79abb30
	github.com/open-cluster-management/multicloud-operators-deployable => github.com/stolostron/multicloud-operators-deployable v1.2.4-0-20211122-24ac963
	github.com/open-cluster-management/multicloud-operators-placementrule => github.com/stolostron/multicloud-operators-placementrule v1.2.4-0-20211122-be034
	github.com/open-cluster-management/multicloud-operators-subscription => github.com/stolostron/multicloud-operators-subscription v1.2.4-0-20211122-7277a37
	github.com/open-cluster-management/multicloud-operators-subscription-release => github.com/stolostron/multicloud-operators-subscription-release v1.2.4-0-20211122-8309641
	github.com/ulikunitz/xz => github.com/ulikunitz/xz v0.5.10
	k8s.io/client-go => k8s.io/client-go v0.23.3
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7
	sigs.k8s.io/application => github.com/kubernetes-sigs/application v0.8.3
)
