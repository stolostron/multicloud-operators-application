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
	pflag "github.com/spf13/pflag"
)

// ControllerRunOptions for the hcm controller.
type ControllerRunOptions struct {
	MetricsAddr          string
	ApplicationCRDFile   string
	LeaderElect          bool
	LeaseDurationSeconds int
	RenewDeadlineSeconds int
	RetryPeriodSeconds   int
}

var options = ControllerRunOptions{
	MetricsAddr:          "",
	ApplicationCRDFile:   "/usr/local/etc/application/crds/app.k8s.io_applications_crd_v1.yaml",
	LeaseDurationSeconds: 137,
	RenewDeadlineSeconds: 107,
	RetryPeriodSeconds:   26,
}

// ProcessFlags parses command line parameters into options
func ProcessFlags() {
	flag := pflag.CommandLine

	// add flags
	flag.StringVar(
		&options.MetricsAddr,
		"metrics-addr",
		options.MetricsAddr,
		"The address the metric endpoint binds to.",
	)

	flag.StringVar(
		&options.ApplicationCRDFile,
		"application-crd-file",
		options.ApplicationCRDFile,
		"Application CRD Yaml File",
	)

	flag.BoolVar(
		&options.LeaderElect,
		"leader-elect",
		false,
		"Enable a leader client to gain leadership before executing the main loop",
	)

	flag.IntVar(
		&options.LeaseDurationSeconds,
		"lease-duration",
		options.LeaseDurationSeconds,
		"The lease duration in seconds.",
	)

	flag.IntVar(
		&options.RenewDeadlineSeconds,
		"renew-deadline",
		options.RenewDeadlineSeconds,
		"The renew deadline in seconds.",
	)

	flag.IntVar(
		&options.RetryPeriodSeconds,
		"retry-period",
		options.RetryPeriodSeconds,
		"The retry period in seconds.",
	)
}
