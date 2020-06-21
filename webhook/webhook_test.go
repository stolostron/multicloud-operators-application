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

package webhook

import (
	"context"
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	mgr "sigs.k8s.io/controller-runtime/pkg/manager"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appv1beta1 "github.com/kubernetes-sigs/application/pkg/apis/app/v1beta1"
)

var _ = Describe("test application validation logic", func() {
	Context("given an exist namespace appliction in a namespace", func() {
		var (
			appkey = types.NamespacedName{Name: "app1", Namespace: "default"}

			appIns = appv1beta1.Application{
				ObjectMeta: metav1.ObjectMeta{
					Name:      appkey.Name,
					Namespace: appkey.Namespace},
				Spec: appv1beta1.ApplicationSpec{
					Selector: &metav1.LabelSelector{
						MatchLabels: map[string]string{"app": "nginx-app-details"},
					},
				},
			}
		)

		BeforeEach(func() {
			// Create the Application object and expect the Reconcile
			Expect(k8sClient.Create(context.TODO(), appIns.DeepCopy())).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(k8sClient.Delete(context.TODO(), &appIns)).Should(Succeed())
		})
	})

	// somehow this test only fail on travis
	// make sure this one runs at the end, otherwise, we might register this
	// webhook before the default one, which cause unexpected results.
	PContext("given a k8s env, it create svc and validating webhook config", func() {
		var (
			lMgr    mgr.Manager
			certDir string
			testNs  string
			caCert  []byte
			err     error
			sstop   chan struct{}
		)

		It("should create a service and ValidatingWebhookConfiguration", func() {
			lMgr, err = mgr.New(testEnv.Config, mgr.Options{MetricsBindAddress: "0"})
			Expect(err).Should(BeNil())

			sstop = make(chan struct{})
			defer close(sstop)
			go func() {
				Expect(lMgr.Start(sstop)).Should(Succeed())
			}()

			certDir = filepath.Join(os.TempDir(), "k8s-webhook-server", "serving-certs")
			testNs = "default"
			os.Setenv("POD_NAMESPACE", testNs)

			caCert, err = GenerateWebhookCerts(certDir)
			Expect(err).Should(BeNil())
			validatorName := "test-validator"
			wbhSvcNm := "app-wbh-svc"
			WireUpWebhookSupplymentryResource(lMgr, stop, wbhSvcNm, validatorName, certDir, caCert)

			ns, err := findEnvVariable(podNamespaceEnvVar)
			Expect(err).Should(BeNil())

			time.Sleep(3 * time.Second)
			wbhSvc := &corev1.Service{}
			svcKey := types.NamespacedName{Name: wbhSvcNm, Namespace: ns}
			Expect(lMgr.GetClient().Get(context.TODO(), svcKey, wbhSvc)).Should(Succeed())
			defer func() {
				Expect(lMgr.GetClient().Delete(context.TODO(), wbhSvc)).Should(Succeed())
			}()

			wbhCfg := &admissionv1.ValidatingWebhookConfiguration{}
			cfgKey := types.NamespacedName{Name: validatorName}
			Expect(lMgr.GetClient().Get(context.TODO(), cfgKey, wbhCfg)).Should(Succeed())

			defer func() {
				Expect(lMgr.GetClient().Delete(context.TODO(), wbhCfg)).Should(Succeed())
			}()
		})
	})
})
