// Copyright 2021 The Kubernetes Authors.
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
	"testing"

	"github.com/onsi/gomega"
	dplv1 "github.com/open-cluster-management/multicloud-operators-deployable/pkg/apis/apps/v1"
	subv1 "github.com/open-cluster-management/multicloud-operators-subscription/pkg/apis/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	appv1beta1 "sigs.k8s.io/application/api/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/event"
)

var (
	chnkey = types.NamespacedName{
		Name:      "test-chn",
		Namespace: "test-chn-namespace",
	}

	oldSubscription = &subv1.Subscription{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Subscription",
			APIVersion: "apps.open-cluster-management.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "subscription1",
			Labels: map[string]string{
				"name": "subscription1",
				"key1": "c1v1",
				"key2": "c1v2",
			},
		},
		Status: subv1.SubscriptionStatus{
			LastUpdateTime: metav1.Now(),
			Reason:         "test",
			Phase:          "Propagated",
		},
		Spec: subv1.SubscriptionSpec{
			Channel: chnkey.String(),
			TimeWindow: &subv1.TimeWindow{
				WindowType: "active",
				Daysofweek: []string{},
				Hours: []subv1.HourRange{
					{Start: "10:00AM", End: "5:00PM"},
				},
			},
		},
	}

	newSubscription = &subv1.Subscription{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Subscription",
			APIVersion: "apps.open-cluster-management.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "subscription1",
			Labels: map[string]string{
				"name": "subscription1",
				"key1": "c1v1",
				"key2": "c1v3",
			},
		},
		Status: subv1.SubscriptionStatus{
			LastUpdateTime: metav1.Now(),
			Reason:         "test",
			Phase:          "Failed",
		},
		Spec: subv1.SubscriptionSpec{
			Channel: chnkey.String(),
			TimeWindow: &subv1.TimeWindow{
				WindowType: "active",
				Daysofweek: []string{},
				Hours: []subv1.HourRange{
					{Start: "09:00AM", End: "5:00PM"},
				},
			},
		},
	}

	oldDeployable = &dplv1.Deployable{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployable",
			APIVersion: "apps.open-cluster-management.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployable1",
			Labels: map[string]string{
				"name": "deployable1",
				"key1": "c1v1",
				"key2": "c1v2",
			},
		},
		Spec: dplv1.DeployableSpec{
			Channels: []string{"test-1"},
		},
	}

	newDeployable = &dplv1.Deployable{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployable",
			APIVersion: "apps.open-cluster-management.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "deployable1",
			Labels: map[string]string{
				"name": "deployable1",
				"key1": "c1v1",
				"key2": "c1v2",
			},
		},
		Spec: dplv1.DeployableSpec{
			Channels: []string{"test-2"},
			Template: &runtime.RawExtension{
				Raw: []byte("ad,123"),
			},
		},
	}
)

func TestUpdateAppInstance(t *testing.T) {
	tests := []struct {
		name     string
		oldApp   *appv1beta1.Application
		newApp   *appv1beta1.Application
		expected bool
	}{
		{
			name:     "empty applications",
			oldApp:   &appv1beta1.Application{},
			newApp:   &appv1beta1.Application{},
			expected: false,
		},
		{
			name: "different subscriptions",
			oldApp: &appv1beta1.Application{ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					"apps.open-cluster-management.io/subscriptions": "default/oldApp",
				},
			}},
			newApp: &appv1beta1.Application{ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					"apps.open-cluster-management.io/subscriptions": "default/newApp",
				},
			}},
			expected: true,
		},
		{
			name: "different deployables",
			oldApp: &appv1beta1.Application{ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					"apps.open-cluster-management.io/deployables": "default/oldApp",
				},
			}},
			newApp: &appv1beta1.Application{ObjectMeta: metav1.ObjectMeta{
				Annotations: map[string]string{
					"apps.open-cluster-management.io/deployables": "default/newApp",
				},
			}},
			expected: true,
		},
	}

	for _, tC := range tests {
		t.Run(tC.name, func(t *testing.T) {
			actual := UpdateAppInstance(tC.oldApp, tC.newApp)
			if actual != tC.expected {
				t.Errorf("UpdateAppInstance expected %v, got %v", tC.expected, actual)
			}
		})
	}
}

func TestPredicate(t *testing.T) {
	g := gomega.NewGomegaWithT(t)

	// Test SubscriptionPredicateFunc
	instance := SubscriptionPredicateFunc

	updateEvt := event.UpdateEvent{
		ObjectOld: oldSubscription,
		ObjectNew: newSubscription,
	}

	ret := instance.Update(updateEvt)
	g.Expect(ret).To(gomega.Equal(true))

	// Test DeployablePredicateFunc
	instance = DeployablePredicateFunc

	updateEvt = event.UpdateEvent{
		ObjectOld: oldDeployable,
		ObjectNew: newDeployable,
	}

	ret = instance.Update(updateEvt)
	g.Expect(ret).To(gomega.Equal(true))
}
