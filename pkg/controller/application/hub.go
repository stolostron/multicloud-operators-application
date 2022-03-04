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
	"strings"

	dplv1 "github.com/open-cluster-management/multicloud-operators-deployable/pkg/apis/apps/v1"
	subv1 "github.com/open-cluster-management/multicloud-operators-subscription/pkg/apis/apps/v1"
	"github.com/stolostron/multicloud-operators-application/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
	appv1beta1 "sigs.k8s.io/application/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *ReconcileApplication) doAppHubReconcile(app *appv1beta1.Application) {
	// allSubs: all subscriptions
	// allDpls: all deployables. The deployables subscribed in the subscriptions are not counted
	// allClusterDplMap: all deployables for each cluster. The deployables subscribed in the subscriptions are counted.
	// All deployables will be required for searching deployed pods
	allSubs, allDpls, allClusterDplMap := r.GetAllNewDeployablesByApplication(app)

	utils.PrintAllClusterDplMap(allClusterDplMap)

	substr := ""
	dplstr := ""

	for _, sub := range allSubs {
		if substr != "" {
			substr += ","
		}

		substr += sub.Namespace + "/" + sub.Name
	}

	for _, dpl := range allDpls {
		if dplstr != "" {
			dplstr += ","
		}

		dplstr += dpl.Namespace + "/" + dpl.Name
	}

	if app.Annotations == nil {
		app.Annotations = make(map[string]string)
	}

	app.Annotations["apps.open-cluster-management.io/subscriptions"] = substr
	app.Annotations["apps.open-cluster-management.io/deployables"] = dplstr
}

// In 2.5, disable setting the part-of label on all subscriptions of the application.
// Subscription controller has updated the same label with different vaule.
// As the result, we see the part-of label is updated to the application name firstly, then it is updated to the subscription name in next cycle.
// Then it is updated to the application name again. so the appsub is updated in a endless loop.
// application crd/controller will deprecate in 2.6.

//GetAllSubscriptionDeployablesByApplication get all subscriptions and their deployables.app.ibm.com objects by a application
func (r *ReconcileApplication) GetAllSubscriptionDeployablesByApplication(app *appv1beta1.Application,
	allClusterDplMap map[string]*utils.DplMap) ([]*subv1.Subscription, error) {
	var allSubs []*subv1.Subscription

	subscriptionList := &subv1.SubscriptionList{}

	listOptions := &client.ListOptions{Namespace: app.Namespace}

	if app.Spec.Selector != nil {
		subSelector, err := utils.ConvertLabels(app.Spec.Selector)
		if err != nil {
			klog.Error("Failed to set label selector of application: ", app.Name, "err: ", err)
		}

		listOptions.LabelSelector = subSelector
	}

	err := r.List(context.TODO(), subscriptionList, listOptions)
	if err != nil {
		klog.Error("Failed to list subscription objects from application namespace ", app.Namespace, " error: ", err)

		if !errors.IsNotFound(err) {
			return nil, nil
		}
	}

	for _, subscription := range subscriptionList.Items {
		allSubs = append(allSubs, subscription.DeepCopy())

		//Check if there is a deployable (its name is "subscrioptionName-deployable") created for deploying the subscription to managed clusters,
		//the deployable status is used for fetching the managed clusters
		subdpl := &dplv1.Deployable{}
		subdplkey := types.NamespacedName{Name: subscription.Name + "-deployable", Namespace: subscription.Namespace}
		err = r.Get(context.TODO(), subdplkey, subdpl)

		if err != nil {
			klog.V(1).Infof("The deployable created for deploying the subscription not found: subdpl: %#v, error: %#v", subdplkey, err)
			continue
		}

		dpls := strings.Split(subscription.Annotations[subv1.AnnotationDeployables], ",")

		for _, dplkey := range dpls {
			strs := strings.Split(dplkey, "/")

			dplSpace := ""
			dplName := ""

			if len(strs) == 2 {
				dplSpace = strs[0]
				dplName = strs[1]
			}

			if dplSpace == "" || dplName == "" {
				continue
			}

			dpl := &dplv1.Deployable{}
			dplkey2 := types.NamespacedName{Name: dplName, Namespace: dplSpace}
			err := r.Get(context.TODO(), dplkey2, dpl)

			if err != nil {
				klog.V(1).Infof("The deployable in the subscription not found: sub: %#v, dpl: %#v, error: %#v", subscription, dplkey2, err)
				continue
			}

			utils.AppendClusterDplMap(*subdpl, *dpl, allClusterDplMap)
		}
	}

	klog.V(1).Infoln("Got all subscriptions in the application: ", app.Name, app.Kind, "|", allSubs)

	return allSubs, nil
}

//GetAllNewDeployablesByApplication get all deployables.app.ibm.com objects by a application
func (r *ReconcileApplication) GetAllNewDeployablesByApplication(
	app *appv1beta1.Application) ([]*subv1.Subscription, []*dplv1.Deployable, map[string]*utils.DplMap) {
	var allSubs []*subv1.Subscription

	var allDpls []*dplv1.Deployable

	allClusterDplMap := make(map[string]*utils.DplMap)

	dplList := &dplv1.DeployableList{}

	dplListOptions := &client.ListOptions{Namespace: app.Namespace}

	if app.Spec.Selector != nil {
		clSelector, err := utils.ConvertLabels(app.Spec.Selector)
		if err != nil {
			klog.Error("Failed to set label selector of application: ", app.Name, "err: ", err)
		}

		dplListOptions.LabelSelector = clSelector
	}

	err := r.List(context.TODO(), dplList, dplListOptions)
	if err != nil {
		klog.Error("Failed to list objects from application namespace ", app.Namespace, " error: ", err)

		if !errors.IsNotFound(err) {
			return nil, nil, nil
		}
	}

	for _, dpl := range dplList.Items {
		if dpl.Annotations != nil && dpl.Annotations[dplv1.AnnotationIsGenerated] == "true" {
			continue
		}

		utils.AppendClusterDplMap(dpl, dpl, allClusterDplMap)

		allDpls = append(allDpls, dpl.DeepCopy())
	}

	allSubs, _ = r.GetAllSubscriptionDeployablesByApplication(app, allClusterDplMap)

	newAllSubs := utils.GetUniqueSubscriptions(allSubs)
	newAllDpls := utils.GetUniqueDeployables(allDpls)
	klog.V(1).Infoln("Got all subscriptions and deployables in the application: ", app.Name, app.Kind, "|", newAllSubs, "|", newAllDpls)

	return newAllSubs, newAllDpls, allClusterDplMap
}

//GetAllApplications get all applications
func (r *ReconcileApplication) GetAllApplications() ([]appv1beta1.Application, error) {
	// find everything with label pointer
	klog.V(1).Infoln("Entering get all Applications")

	var applist *appv1beta1.ApplicationList

	listOptions := &client.ListOptions{}

	err := r.List(context.TODO(), applist, listOptions)
	if err != nil {
		klog.Error("Failed to list all application: ", err)
	}

	klog.V(1).Infoln("Get all Applications: ", applist.Items, " len: ", len(applist.Items), " error: ", err)

	return applist.Items, nil
}
