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

package exec

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/open-cluster-management/multicloud-operators-application/pkg/apis"
	"github.com/open-cluster-management/multicloud-operators-application/pkg/controller"
	"github.com/open-cluster-management/multicloud-operators-application/utils"
	appWebhook "github.com/open-cluster-management/multicloud-operators-application/webhook"

	appapis "github.com/kubernetes-sigs/application/pkg/apis"

	dplapis "github.com/open-cluster-management/multicloud-operators-deployable/pkg/apis"
	dplv1 "github.com/open-cluster-management/multicloud-operators-deployable/pkg/apis/apps/v1"

	subapis "github.com/open-cluster-management/multicloud-operators-subscription/pkg/apis"
	subv1 "github.com/open-cluster-management/multicloud-operators-subscription/pkg/apis/apps/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/client-go/rest"
	"k8s.io/klog"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
)

// Change below variables to serve metrics on different host or port.
var (
	metricsHost             = "0.0.0.0"
	metricsPort         int = 8386
	operatorMetricsPort int = 8689
)

// RunManager starts the actual manager
func RunManager() {
	// Get a config to talk to the apiserver
	cfg, err := config.GetConfig()
	if err != nil {
		klog.Error(err, "")
		os.Exit(1)
	}

	runtimeClient, err := client.New(cfg, client.Options{})
	if err != nil {
		klog.Infof("Error building runtime clientset: %s", err)
		os.Exit(1)
	}

	// Register application CRD into hub kubernetes cluster
	err = utils.CheckAndInstallCRD(cfg, options.ApplicationCRDFile)
	if err != nil {
		klog.Infof("unable to install placementrule crd in hub: %s", err)
		os.Exit(1)
	}

	enableLeaderElection := false

	if _, err := rest.InClusterConfig(); err == nil {
		klog.Info("LeaderElection enabled as running in a cluster")

		enableLeaderElection = true
	} else {
		klog.Info("LeaderElection disabled as not running in a cluster")
	}

	// Create a new Cmd to provide shared dependencies and start components
	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		MetricsBindAddress:      fmt.Sprintf("%s:%d", metricsHost, metricsPort),
		Port:                    operatorMetricsPort,
		LeaderElection:          enableLeaderElection,
		LeaderElectionID:        "multicloud-operators-application-leader.open-cluster-management.io",
		LeaderElectionNamespace: "kube-system",
	})
	if err != nil {
		klog.Error(err, "")
		os.Exit(1)
	}

	klog.Info("Registering Components.")

	// Setup Scheme for all resources
	if err := apis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Error(err, "")
		os.Exit(1)
	}

	//append deployables.apps.open-cluster-management.io and subscriptions.apps.open-cluster-management to scheme
	if err = dplapis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Error("unable add deployables.apps.open-cluster-management.io APIs to scheme: ", err)
		os.Exit(1)
	}

	if err = subapis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Error("unable add subscriptions.apps.open-cluster-management.io APIs to scheme: ", err)
		os.Exit(1)
	}

	//append application api to scheme
	if err = appapis.AddToScheme(mgr.GetScheme()); err != nil {
		klog.Error("unable add mcm APIs to scheme: ", err)
		os.Exit(1)
	}

	dpllist := &dplv1.DeployableList{}
	err = runtimeClient.List(context.TODO(), dpllist, &client.ListOptions{})

	if err != nil && !errors.IsNotFound(err) {
		klog.Fatal("Deployable kind is not ready in api server, exit and retry later")
		os.Exit(1)
	}

	sublist := &subv1.SubscriptionList{}
	err = runtimeClient.List(context.TODO(), sublist, &client.ListOptions{})

	if err != nil && !errors.IsNotFound(err) {
		klog.Fatal("Subscription kind is not ready in api server, exit and retry later")
		os.Exit(1)
	}

	// Setup all Controllers
	if err := controller.AddToManager(mgr); err != nil {
		klog.Error(err, "")
		os.Exit(1)
	}

	sig := signals.SetupSignalHandler()

	// Setup webhooks
	klog.Info("setting up webhook server")

	clt, err := client.New(ctrl.GetConfigOrDie(), client.Options{})
	if err != nil {
		klog.Errorf("failed to create a client for webhook to get CA cert secret, err %v", err)
		os.Exit(1)
	}

	hookServer := mgr.GetWebhookServer()
	certDir := filepath.Join(os.TempDir(), "k8s-webhook-server", "application-serving-certs")

	caCert, err := appWebhook.WireUpWebhook(clt, mgr, hookServer, certDir)
	if err != nil {
		klog.Error(err, "failed to wire up webhook")
		os.Exit(1)
	}

	go appWebhook.WireUpWebhookSupplymentryResource(mgr, sig, appWebhook.WebhookServiceName,
		appWebhook.WebhookValidatorName, certDir, caCert)

	klog.Info("Starting the Cmd.")

	// Start the Cmd
	if err := mgr.Start(sig); err != nil {
		klog.Error(err, "Manager exited non-zero")
		os.Exit(1)
	}
}
