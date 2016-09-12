/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package healthchecks

import (
	"fmt"

	compute "google.golang.org/api/compute/v1"
	"k8s.io/contrib/ingress/controllers/gce/utils"
)

// NewFakeHealthChecks returns a new FakeHealthChecks.
func NewFakeHealthChecks() *FakeHealthChecks {
	return &FakeHealthChecks{hc: []*compute.HttpHealthCheck{}}
}

// FakeHealthCheckGetter implements the healthCheckGetter interface for tests.
type FakeHealthCheckGetter struct {
	DefaultHealthCheck *compute.HttpHealthCheck
}

// HealthCheck returns the health check for the given port. If a health check
// isn't stored under the DefaultHealthCheck member, it constructs one.
func (h *FakeHealthCheckGetter) HealthCheck(port int64) (*compute.HttpHealthCheck, error) {
	if h.DefaultHealthCheck == nil {
		return utils.DefaultHealthCheckTemplate(port), nil
	}
	return h.DefaultHealthCheck, nil
}

// FakeHealthChecks fakes out health checks.
type FakeHealthChecks struct {
	hc []*compute.HttpHealthCheck
}

// CreateHttpHealthCheck fakes out http health check creation.
func (f *FakeHealthChecks) CreateHttpHealthCheck(hc *compute.HttpHealthCheck) error {
	f.hc = append(f.hc, hc)
	return nil
}

// GetHttpHealthCheck fakes out getting a http health check from the cloud.
func (f *FakeHealthChecks) GetHttpHealthCheck(name string) (*compute.HttpHealthCheck, error) {
	for _, h := range f.hc {
		if h.Name == name {
			return h, nil
		}
	}
	return nil, fmt.Errorf("Health check %v not found.", name)
}

// DeleteHttpHealthCheck fakes out deleting a http health check.
func (f *FakeHealthChecks) DeleteHttpHealthCheck(name string) error {
	healthChecks := []*compute.HttpHealthCheck{}
	exists := false
	for _, h := range f.hc {
		if h.Name == name {
			exists = true
			continue
		}
		healthChecks = append(healthChecks, h)
	}
	if !exists {
		return fmt.Errorf("Failed to find health check %v", name)
	}
	f.hc = healthChecks
	return nil
}

// UpdateHttpHealthCheck sends the given health check as an update.
func (f *FakeHealthChecks) UpdateHttpHealthCheck(hc *compute.HttpHealthCheck) error {
	healthChecks := []*compute.HttpHealthCheck{}
	found := false
	for _, h := range f.hc {
		if h.Name == hc.Name {
			healthChecks = append(healthChecks, hc)
			found = true
		} else {
			healthChecks = append(healthChecks, h)
		}
	}
	if !found {
		return fmt.Errorf("Cannot update a non-existent health check %v", hc.Name)
	}
	f.hc = healthChecks
	return nil
}
