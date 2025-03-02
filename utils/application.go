package utils

import (
	"encoding/json"

	"reflect"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	dplv1 "github.com/stolostron/multicloud-operators-application/pkg/apis/deployable/v1"
	subv1 "open-cluster-management.io/multicloud-operators-subscription/pkg/apis/apps/v1"
	appv1beta1 "sigs.k8s.io/application/api/v1beta1"
)

// DeployablePredicateFunc defines predicate function for deployable watch in deployable controller
var DeployablePredicateFunc = predicate.Funcs{
	UpdateFunc: func(e event.UpdateEvent) bool {
		newdpl := e.ObjectNew.(*dplv1.Deployable)
		olddpl := e.ObjectOld.(*dplv1.Deployable)

		if len(newdpl.GetFinalizers()) > 0 {
			return true
		}

		if !reflect.DeepEqual(newdpl.GetAnnotations(), olddpl.GetAnnotations()) {
			return true
		}

		if !reflect.DeepEqual(newdpl.GetLabels(), olddpl.GetLabels()) {
			return true
		}

		oldtmpl := &unstructured.Unstructured{}
		newtmpl := &unstructured.Unstructured{}

		if olddpl.Spec.Template == nil || olddpl.Spec.Template.Raw == nil {
			return true
		}
		err := json.Unmarshal(olddpl.Spec.Template.Raw, oldtmpl)
		if err != nil {
			return true
		}
		if newdpl.Spec.Template.Raw == nil {
			return true
		}
		err = json.Unmarshal(newdpl.Spec.Template.Raw, newtmpl)
		if err != nil {
			return true
		}

		if !reflect.DeepEqual(newtmpl, oldtmpl) {
			return true
		}

		olddpl.Spec.Template = newdpl.Spec.Template.DeepCopy()
		return !reflect.DeepEqual(olddpl.Spec, newdpl.Spec)
	},
}

// SubscriptionPredicateFunc filters status update
var SubscriptionPredicateFunc = predicate.Funcs{
	UpdateFunc: func(e event.UpdateEvent) bool {
		subOld := e.ObjectOld.(*subv1.Subscription)
		subNew := e.ObjectNew.(*subv1.Subscription)

		// need to process delete with finalizers
		if len(subNew.GetFinalizers()) > 0 {
			return true
		}

		// we care label change, pass it down
		if !reflect.DeepEqual(subOld.GetLabels(), subNew.GetLabels()) {
			return true
		}

		// we care annotation change. pass it down
		if !reflect.DeepEqual(subOld.GetAnnotations(), subNew.GetAnnotations()) {
			return true
		}

		// we care spec for sure
		if !reflect.DeepEqual(subOld.Spec, subNew.Spec) {
			return true
		}

		// do we care phase change?
		if subNew.Status.Phase == "" || subNew.Status.Phase != subOld.Status.Phase {
			klog.V(5).Info("We care phase..", subNew.Status.Phase, " vs ", subOld.Status.Phase)
			return true
		}

		klog.V(1).Info("Something we don't care changed")
		return false
	},
}

func UpdateAppInstance(oldApp, newApp *appv1beta1.Application) bool {
	//check dpl and subscription list annotations
	oldAppAnno := oldApp.GetAnnotations()
	if oldAppAnno == nil {
		oldAppAnno = make(map[string]string)
	}

	newAppAnno := newApp.GetAnnotations()
	if newAppAnno == nil {
		newAppAnno = make(map[string]string)
	}

	if !reflect.DeepEqual(oldAppAnno["apps.open-cluster-management.io/subscriptions"], newAppAnno["apps.open-cluster-management.io/subscriptions"]) {
		return true
	}

	if !reflect.DeepEqual(oldAppAnno["apps.open-cluster-management.io/deployables"], newAppAnno["apps.open-cluster-management.io/deployables"]) {
		return true
	}

	return false
}
