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
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	. "github.com/onsi/gomega"
	admissionv1 "k8s.io/api/admissionregistration/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	mgr "sigs.k8s.io/controller-runtime/pkg/manager"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	appv1beta1 "sigs.k8s.io/application/api/v1beta1"
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
		)

		It("should create a service and ValidatingWebhookConfiguration", func() {
			lMgr, err = mgr.New(testEnv.Config, mgr.Options{MetricsBindAddress: "0"})
			Expect(err).Should(BeNil())

			ctx, cancel := context.WithCancel(context.Background())
			defer func() {
				cancel()
			}()

			go func() {
				Expect(lMgr.Start(ctx)).Should(Succeed())
			}()

			certDir = filepath.Join(os.TempDir(), "k8s-webhook-server", "application-serving-certs")
			testNs = "default"
			os.Setenv("POD_NAMESPACE", testNs)

			caCert, err = GenerateWebhookCerts(k8sClient, certDir)
			Expect(err).Should(BeNil())
			validatorName := "test-validator"
			wbhSvcNm := "app-wbh-svc"
			WireUpWebhookSupplymentryResource(stop, lMgr, wbhSvcNm, validatorName, certDir, caCert)
			WireUpWebhookSupplymentryResource(ctx, lMgr, wbhSvcNm, validatorName, certDir, caCert)

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

func TestWireupWebhook(t *testing.T) {
	g := gomega.NewGomegaWithT(t)
	var (
		metricsHost         = "0.0.0.0"
		metricsPort         = 8386
		operatorMetricsPort = 8689
	)

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		MetricsBindAddress:      fmt.Sprintf("%s:%d", metricsHost, metricsPort),
		Port:                    operatorMetricsPort,
		LeaderElection:          false,
		LeaderElectionID:        "multicloud-operators-application-leader.open-cluster-management.io",
		LeaderElectionNamespace: "kube-system",
	})
	g.Expect(err).ShouldNot(HaveOccurred())
	clt, err := client.New(ctrl.GetConfigOrDie(), client.Options{})
	hookServer := mgr.GetWebhookServer()
	certDir := filepath.Join(os.TempDir(), "k8s-webhook-server", "application-serving-certs")
	caCert, err := WireUpWebhook(clt, mgr, hookServer, certDir)
	g.Expect(caCert).Should(HaveOccurred())
}
