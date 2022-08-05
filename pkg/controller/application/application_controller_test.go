// Copyright 2019 The Kubernetes Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package application

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	dplv1 "github.com/open-cluster-management/multicloud-operators-deployable/pkg/apis/apps/v1"
	subv1 "github.com/open-cluster-management/multicloud-operators-subscription/pkg/apis/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appv1beta1 "sigs.k8s.io/application/api/v1beta1"
)

var (
	applicationNS = "kube-system"
	chnkey        = types.NamespacedName{
		Name:      "test-chn",
		Namespace: "test-chn-namespace",
	}
)

func TestReconcile(t *testing.T) {
	defer klog.Flush()

	g := gomega.NewGomegaWithT(t)

	// Setup the Manager and Controller.  Wrap the Controller Reconcile function so it writes each request to a
	// channel when it is finished.

	t.Log("Create manager")

	mgr, err := manager.New(cfg, manager.Options{
		MetricsBindAddress: "0",
		LeaderElection:     false,
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	t.Log("Setup test reconcile")
	g.Expect(Add(mgr)).NotTo(gomega.HaveOccurred())

	ctx, cancel := context.WithTimeout(context.TODO(), 5*time.Minute)
	mgrStopped := StartTestManager(ctx, mgr, g)

	defer func() {
		cancel()
		mgrStopped.Wait()
	}()

	deployableName := "example-subscription-deployable"
	deployableInstance := &dplv1.Deployable{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployable",
			APIVersion: "apps.open-cluster-management.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      deployableName,
			Namespace: applicationNS,
		},
		Spec: dplv1.DeployableSpec{
			Template: &runtime.RawExtension{},
		},
	}
	deployableInstance.Spec.Template.Raw, _ = json.Marshal(deployableInstance)

	applicationName := "example-application"
	applicationKey := types.NamespacedName{
		Name:      applicationName,
		Namespace: applicationNS,
	}
	instance := &appv1beta1.Application{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Application",
			APIVersion: "app.k8s.io/v1beta1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      applicationName,
			Namespace: applicationNS,
		},
		Spec: appv1beta1.ApplicationSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": "nginx-app-details"},
			},
		},
	}

	subscriptionName := "example-subscription"
	subscriptionKey := types.NamespacedName{
		Name:      subscriptionName,
		Namespace: "kube-system",
	}
	subInstance := &subv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      subscriptionName,
			Namespace: "kube-system",
		},
		Spec: subv1.SubscriptionSpec{
			Channel: chnkey.String(),
		},
	}

	c := mgr.GetClient()

	err = c.Create(context.TODO(), deployableInstance)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = c.Create(context.TODO(), instance)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = c.Create(context.TODO(), subInstance)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	time.Sleep(4 * time.Second)

	instanceResp := &appv1beta1.Application{}
	err = c.Get(context.TODO(), applicationKey, instanceResp)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	subscriptionResp := &subv1.Subscription{}
	err = c.Get(context.TODO(), subscriptionKey, subscriptionResp)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	slist := &subv1.SubscriptionList{}
	err = c.List(context.TODO(), slist, &client.ListOptions{})
	g.Expect(!errors.IsNotFound(err)).To(gomega.BeTrue())

	instanceResp1 := &subv1.Subscription{}
	err = c.Get(context.TODO(), chnkey, instanceResp1)
	g.Expect(errors.IsNotFound(err)).To(gomega.BeTrue())
}
