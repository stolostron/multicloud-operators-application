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

	dplv1 "github.com/open-cluster-management/multicloud-operators-deployable/pkg/apis/apps/v1"
	subv1 "github.com/open-cluster-management/multicloud-operators-subscription/pkg/apis/apps/v1"
	"github.com/stolostron/multicloud-operators-application/utils"

	"strings"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	appv1beta1 "sigs.k8s.io/application/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Add creates a new Application Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	erecorder, _ := utils.NewEventRecorder(mgr.GetConfig(), mgr.GetScheme())

	return &ReconcileApplication{
		Client:        mgr.GetClient(),
		scheme:        mgr.GetScheme(),
		eventRecorder: erecorder,
	}
}

type deployableMapper struct {
	client.Client
}

func (mapper *deployableMapper) Map(obj client.Object) []reconcile.Request {
	//enqueue all applications under these namespaces including the deployable namespace plus all of subscription namespaces related to the deployable
	dplName := obj.GetName()
	dplNamespace := obj.GetNamespace()
	klog.V(1).Info("In deployable Mapper:", dplName, "/", dplNamespace)

	nsmap := make(map[string]bool)
	nsmap[dplNamespace] = true

	var requests []reconcile.Request

	subscriptionList := &subv1.SubscriptionList{}
	listOptions := &client.ListOptions{}
	err := mapper.List(context.TODO(), subscriptionList, listOptions)

	if err != nil {
		klog.Error("Failed to list all subscription objects. ", "error: ", err)
		return requests
	}

	for _, subscription := range subscriptionList.Items {
		strs := strings.Split(subscription.Spec.Channel, "/")
		if len(strs) == 2 {
			subChannelns := strs[0]
			if subChannelns == dplNamespace {
				nsmap[subscription.Namespace] = true
			}
		}
	}

	applicationList := &appv1beta1.ApplicationList{}
	err = mapper.List(context.TODO(), applicationList, listOptions)

	if err != nil {
		klog.Error("Failed to list all subscription objects. ", "error: ", err)
		return requests
	}

	for _, app := range applicationList.Items {
		if nsmap[app.GetNamespace()] {
			objkey := types.NamespacedName{
				Name:      app.GetName(),
				Namespace: app.GetNamespace(),
			}

			requests = append(requests, reconcile.Request{NamespacedName: objkey})
		}
	}

	return requests
}

type subscriptionMapper struct {
	client.Client
}

func (mapper *subscriptionMapper) Map(obj client.Object) []reconcile.Request {
	//enqueue all applications under the subscription namespace
	subName := obj.GetName()
	subNamespace := obj.GetNamespace()
	klog.V(1).Info("In subscription Mapper:", subName, "/", subNamespace)

	var requests []reconcile.Request

	applicationList := &appv1beta1.ApplicationList{}
	listOptions := &client.ListOptions{Namespace: subNamespace}
	err := mapper.List(context.TODO(), applicationList, listOptions)

	if err != nil {
		klog.Error("Failed to list all application objects. ", "error: ", err)
		return requests
	}

	for _, app := range applicationList.Items {
		objkey := types.NamespacedName{
			Name:      app.GetName(),
			Namespace: app.GetNamespace(),
		}

		requests = append(requests, reconcile.Request{NamespacedName: objkey})
	}

	return requests
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("application-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Application
	err = c.Watch(&source.Kind{Type: &appv1beta1.Application{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to Deployable
	dmapper := &deployableMapper{mgr.GetClient()}

	err = c.Watch(
		&source.Kind{Type: &dplv1.Deployable{}},
		handler.EnqueueRequestsFromMapFunc(dmapper.Map),
		utils.DeployablePredicateFunc)
	if err != nil {
		return err
	}

	// Watch for changes to Subscription
	smapper := &subscriptionMapper{mgr.GetClient()}

	err = c.Watch(
		&source.Kind{Type: &subv1.Subscription{}},
		handler.EnqueueRequestsFromMapFunc(smapper.Map),
		utils.SubscriptionPredicateFunc)
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileApplication implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileApplication{}

// ReconcileApplication reconciles a Application object
type ReconcileApplication struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client.Client
	scheme        *runtime.Scheme
	eventRecorder *utils.EventRecorder
}

// Reconcile reads that state of the cluster for a Application object and makes changes based on the state read
// and what is in the Application.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileApplication) Reconcile(ctx context.Context, request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Deployable instance
	instance := &appv1beta1.Application{}
	err := r.Get(ctx, request.NamespacedName, instance)
	klog.Info("Reconciling Application:", request.NamespacedName, " with Get err:", err)

	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// validate all deployables, remove the deployables whose hosting deployables are gone
			klog.Info("Reconciling - finished.", request.NamespacedName, " with Get err:", err)

			return reconcile.Result{}, err
		}
		// Error reading the object - requeue the request.
		klog.Info("Reconciling - finished.", request.NamespacedName, " with Get err:", err)

		return reconcile.Result{}, err
	}

	oldInstance := instance.DeepCopy()

	r.doAppHubReconcile(instance)

	result := reconcile.Result{}

	if utils.UpdateAppInstance(oldInstance, instance) {
		klog.V(1).Infoln("Update app annotation", instance.Annotations)

		addtionalMsg := "The app annotations updated. App:" + instance.Namespace + "/" + instance.Name
		r.eventRecorder.RecordEvent(instance, "Update", addtionalMsg, nil)

		err = r.Update(ctx, instance)
		if err != nil {
			klog.Error("Error returned when updating application :", err, "instance:", instance.GetNamespace()+"/"+instance.GetName())
			return reconcile.Result{}, err
		}
	}

	return result, nil
}
