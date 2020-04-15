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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	appv1beta1 "github.com/kubernetes-sigs/application/pkg/apis/app/v1beta1"
	dplv1 "github.com/open-cluster-management/multicloud-operators-deployable/pkg/apis/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	applicationNS = "kube-system"
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

	stopMgr, mgrStopped := StartTestManager(mgr, g)

	defer func() {
		close(stopMgr)
		mgrStopped.Wait()
	}()

	deployableName := "example-deployable"
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

	c := mgr.GetClient()

	err = c.Create(context.TODO(), deployableInstance)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	err = c.Create(context.TODO(), instance)
	g.Expect(err).NotTo(gomega.HaveOccurred())

	time.Sleep(4 * time.Second)

	instanceResp := &appv1beta1.Application{}
	err = c.Get(context.TODO(), applicationKey, instanceResp)
	g.Expect(err).NotTo(gomega.HaveOccurred())
}
