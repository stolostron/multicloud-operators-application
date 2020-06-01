// Copyright 2020 The Kubernetes Authors.
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

	deployablev1 "github.com/open-cluster-management/multicloud-operators-deployable/pkg/apis/apps/v1"
	subv1 "github.com/open-cluster-management/multicloud-operators-subscription/pkg/apis/apps/v1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	dplname = "example-configmap"
	dplns   = "default"
)

var (
	payload = &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			Kind: "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "payload",
		},
	}
)

func TestGetUniqueDeployables(t *testing.T) {
	instance := &deployablev1.Deployable{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dplname,
			Namespace: dplns,
		},
		Spec: deployablev1.DeployableSpec{
			Template: &runtime.RawExtension{
				Object: payload,
			},
		},
	}

	var dplArray []*deployablev1.Deployable
	dplArray = append(dplArray, instance)

	dpl := GetUniqueDeployables(dplArray)
	assert.Equal(t, dpl[0].GetName(), instance.GetName())
}

func TestGetUniqueSubscriptions(t *testing.T) {
	instance := &subv1.Subscription{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-sub",
			Namespace: "test-sub-namespace",
		},
	}

	var subArray []*subv1.Subscription
	subArray = append(subArray, instance)

	sub := GetUniqueSubscriptions(subArray)
	assert.Equal(t, sub[0].GetName(), instance.GetName())
}
