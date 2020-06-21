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
	"encoding/json"
	"fmt"
	"net/http"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	appv1beta1 "github.com/kubernetes-sigs/application/pkg/apis/app/v1beta1"
)

type AppValidator struct {
	client.Client
	decoder *admission.Decoder
}

// AppValidator denys a application creat/update if the application had bad input like this
// the `values` is defined as a string rather than a list.
//  selector:
//    matchExpressions:
//    - key: app
//	    operator: In
//	    values: val-app-1

func (v *AppValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	log.Info("entry webhook handle")
	defer log.Info("exit webhook handle")

	app := &appv1beta1.Application{}

	err := v.decoder.Decode(req, app)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	appJSON, err := json.Marshal(app)
	if err != nil {
		return admission.Denied(fmt.Sprint("convert to JSON failed: ", err))
	}

	newApp := &appv1beta1.Application{}
	err = json.Unmarshal(appJSON, newApp)

	if err != nil {
		return admission.Denied(fmt.Sprint("Invalid application object: ", err))
	}

	return admission.Allowed("")
}

// AppValidator implements admission.DecoderInjector.
// A decoder will be automatically injected.

// InjectDecoder injects the decoder.
func (v *AppValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
