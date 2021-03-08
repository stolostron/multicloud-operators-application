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

package utils

import (
	"encoding/json"

	dplv1alpha1 "github.com/open-cluster-management/multicloud-operators-deployable/pkg/apis/apps/v1"
	subv1alpha1 "github.com/open-cluster-management/multicloud-operators-subscription/pkg/apis/apps/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog"
)

const (
	ResourceLabelName = "name"
)

//DplMap store all dpl names for a cluster [dplName].
type DplMap struct {
	DplResourceMap map[string]*dplv1alpha1.Deployable
}

//GetUniqueDeployables get unique deployable array.
func GetUniqueDeployables(allDpls []*dplv1alpha1.Deployable) []*dplv1alpha1.Deployable {
	dplmap := make(map[string]*dplv1alpha1.Deployable)

	for _, dpl := range allDpls {
		dplkey := types.NamespacedName{Name: dpl.Name, Namespace: dpl.Namespace}.String()
		dplmap[dplkey] = dpl.DeepCopy()
	}

	var newdpls []*dplv1alpha1.Deployable
	for _, newdpl := range dplmap {
		newdpls = append(newdpls, newdpl.DeepCopy())
	}

	return newdpls
}

//GetUniqueSubscriptions get unique subscription array.
func GetUniqueSubscriptions(allSubs []*subv1alpha1.Subscription) []*subv1alpha1.Subscription {
	submap := make(map[string]*subv1alpha1.Subscription)

	for _, sub := range allSubs {
		subkey := types.NamespacedName{Name: sub.Name, Namespace: sub.Namespace}.String()
		submap[subkey] = sub.DeepCopy()
	}

	var newsubs []*subv1alpha1.Subscription
	for _, newsub := range submap {
		newsubs = append(newsubs, newsub.DeepCopy())
	}

	return newsubs
}

//PrintAllClusterDplMap print all cluster deployable map.
func PrintAllClusterDplMap(allClusterDplMap map[string]*DplMap) {
	for cluster, dplmap := range allClusterDplMap {
		for dplname, dpl := range dplmap.DplResourceMap {
			template := &unstructured.Unstructured{}
			templateKind := ""

			if dpl.Spec.Template != nil {
				err := json.Unmarshal(dpl.Spec.Template.Raw, template)
				if err == nil {
					templateKind = template.GetKind()
				}
			}

			klog.V(1).Infof("cluster: %#v, dpl: %#v, dpl template kind: %#v", cluster, dplname, templateKind)
		}
	}
}

//AppendClusterDplMap append dpl and its deployed cluster to allClusterDplMap.
func AppendClusterDplMap(statusdpl dplv1alpha1.Deployable, dpl dplv1alpha1.Deployable, allClusterDplMap map[string]*DplMap) {
	if statusdpl.Status.Phase == "Propagated" || statusdpl.Status.Phase == "Deployed" {
		for cluster := range statusdpl.Status.PropagatedStatus {
			dplmap, ok := allClusterDplMap[cluster]

			if !ok {
				// register new dpl name
				dplmap = &DplMap{
					DplResourceMap: make(map[string]*dplv1alpha1.Deployable),
				}
			}

			dplmap.DplResourceMap[dpl.Name] = dpl.DeepCopy()
			allClusterDplMap[cluster] = dplmap
		}
	}
}
