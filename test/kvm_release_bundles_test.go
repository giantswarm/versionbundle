package test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/giantswarm/versionbundle"
	"github.com/go-resty/resty"
)

func TestKVMReleaseBundles(t *testing.T) {
	versionBundleEndpointHandlerFuncs := []func(w http.ResponseWriter, r *http.Request){
		certOperatorVersionBundlesEndpoint(t),
		clusterOperatorVersionBundlesEndpoint(t),
		flannelOperatorVersionBundlesEndpoint(t),
		kvmOperatorVersionBundlesEndpoint(t),
	}

	var endpoints []*url.URL
	{
		for _, hf := range versionBundleEndpointHandlerFuncs {
			ts := httptest.NewServer(http.HandlerFunc(hf))
			defer ts.Close()
			u, err := url.Parse(ts.URL)
			if err != nil {
				t.Fatalf("expected %#v got %#v", nil, err)
			}
			endpoints = append(endpoints, u)
		}
	}

	var err error

	var collector *versionbundle.Collector
	{
		c := versionbundle.CollectorConfig{
			// Only allow provider independent or KVM provider bundles
			FilterFunc: func(b versionbundle.Bundle) bool { return b.Provider == "" || b.Provider == "kvm" },
			Logger:     microloggertest.New(),
			RestClient: resty.New(),

			Endpoints: endpoints,
		}

		collector, err = versionbundle.NewCollector(c)
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
	}

	var aggregator *versionbundle.Aggregator
	{
		c := versionbundle.AggregatorConfig{
			Logger: microloggertest.New(),
		}

		aggregator, err = versionbundle.NewAggregator(c)
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
	}

	err = collector.Collect(context.TODO())
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}

	aggregated, err := aggregator.Aggregate(collector.Bundles())
	if err != nil {
		t.Fatalf("expected %#v got %#v", nil, err)
	}

	var releases []versionbundle.Release
	for _, a := range aggregated {
		c := versionbundle.DefaultReleaseConfig()

		c.Bundles = a

		r, err := versionbundle.NewRelease(c)
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}

		releases = append(releases, r)
	}

	//
	// Validation
	//
	expectedReleaseCount := 22
	expectedMinReleaseVersion := "0.4.0"
	expectedMaxReleaseVersion := "2.5.4"
	expectedBundleCountInEachRelease := 4

	if len(releases) != expectedReleaseCount {
		t.Fatalf("expected %d releases got %d", expectedReleaseCount, len(releases))
	}

	minReleaseVersion := releases[0].Version()
	maxReleaseVersion := releases[0].Version()
	for _, r := range releases {
		// Find min & max versions in releases.
		if r.Version() < minReleaseVersion {
			minReleaseVersion = r.Version()
		}

		if r.Version() > maxReleaseVersion {
			maxReleaseVersion = r.Version()
		}

		bb := r.Bundles()

		// Verify bundle count in each release.
		if len(bb) != expectedBundleCountInEachRelease {
			t.Fatalf("expected bundle count in each release %d, in release %s got %d", expectedBundleCountInEachRelease, r.Version(), len(bb))
		}

		// Verify Bundle.Provider is provider independent or for KVM
		for _, b := range bb {
			if b.Provider != "" && b.Provider != "kvm" {
				t.Fatalf("expected Bundle.Provider to be empty or kvm, got %s", b.Provider)
			}
		}

		// Verify there are no duplicate bundles in release.
		for i := 1; i < len(bb); i++ {
			if bb[i-1].Name == bb[i].Name {
				t.Fatalf("in release %s there is duplicate bundles for name %s", r.Version(), bb[i].Name)
			}
		}
	}

	if minReleaseVersion != expectedMinReleaseVersion {
		t.Fatalf("expected minReleaseVersion == %s, got %s", expectedMinReleaseVersion, minReleaseVersion)
	}

	if maxReleaseVersion != expectedMaxReleaseVersion {
		t.Fatalf("expected maxReleaseVersion == %s, got %s", expectedMaxReleaseVersion, maxReleaseVersion)
	}
}

func certOperatorVersionBundlesEndpoint(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cr := versionbundle.CollectorEndpointResponse{
			VersionBundles: []versionbundle.Bundle{
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "vault",
							Description: "Vault version updated.",
							Kind:        versionbundle.KindChanged,
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "vault",
							Version: "0.7.3",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   false,
					Name:         "cert-operator",
					Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
					Version:      "0.1.0",
					WIP:          false,
				},
			},
		}
		b, err := json.Marshal(cr)
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
		_, err = io.WriteString(w, string(b))
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
	}
}

func clusterOperatorVersionBundlesEndpoint(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cr := versionbundle.CollectorEndpointResponse{
			VersionBundles: []versionbundle.Bundle{
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "Cluster Operator",
							Description: "Initial version for AWS",
							Kind:        "added",
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "aws-operator",
							Version: "1.0.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   false,
					Name:         "cluster-operator",
					Provider:     "aws",
					Time:         time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
					Version:      "0.1.0",
					WIP:          true,
				},
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "Cluster Operator",
							Description: "Initial version for Azure",
							Kind:        "added",
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "azure-operator",
							Version: "1.0.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   false,
					Name:         "cluster-operator",
					Provider:     "azure",
					Time:         time.Date(2018, time.March, 28, 7, 30, 0, 0, time.UTC),
					Version:      "0.1.0",
					WIP:          true,
				},
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "Cluster Operator",
							Description: "Initial version for KVM",
							Kind:        "added",
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "kvm-operator",
							Version: "1.0.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   false,
					Name:         "cluster-operator",
					Provider:     "kvm",
					Time:         time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
					Version:      "0.1.0",
					WIP:          true,
				},
			},
		}
		b, err := json.Marshal(cr)
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
		_, err = io.WriteString(w, string(b))
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
	}
}

func flannelOperatorVersionBundlesEndpoint(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cr := versionbundle.CollectorEndpointResponse{
			VersionBundles: []versionbundle.Bundle{
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "flannel",
							Description: "Flannel version updated.",
							Kind:        "changed",
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "flannel",
							Version: "0.9.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   false,
					Name:         "flannel-operator",
					Time:         time.Date(2017, time.October, 27, 16, 21, 0, 0, time.UTC),
					Version:      "0.1.0",
					WIP:          false,
				},
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "flannel",
							Description: "Flannel version updated.",
							Kind:        "changed",
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "flannel",
							Version: "0.10.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   false,
					Name:         "flannel-operator",
					Time:         time.Date(2018, time.March, 16, 9, 15, 0, 0, time.UTC),
					Version:      "0.2.0",
					WIP:          false,
				},
			},
		}
		b, err := json.Marshal(cr)
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
		_, err = io.WriteString(w, string(b))
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
	}
}

func kvmOperatorVersionBundlesEndpoint(t *testing.T) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cr := versionbundle.CollectorEndpointResponse{
			VersionBundles: []versionbundle.Bundle{
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "docker",
							Description: "Docker version updated.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "etcd",
							Description: "Etcd version updated.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "kubedns",
							Description: "KubeDNS version updated.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version updated.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "nginx-ingress-controller",
							Description: "Nginx-ingress-controller version updated.",
							Kind:        versionbundle.KindChanged,
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "calico",
							Version: "2.6.2",
						},
						{
							Name:    "docker",
							Version: "1.12.6",
						},
						{
							Name:    "etcd",
							Version: "3.2.7",
						},
						{
							Name:    "kubedns",
							Version: "1.14.5",
						},
						{
							Name:    "kubernetes",
							Version: "1.8.1",
						},
						{
							Name:    "nginx-ingress-controller",
							Version: "0.9.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2017, time.October, 26, 16, 38, 0, 0, time.UTC),
					Version:      "0.1.0",
					WIP:          false,
				},
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "kubernetes",
							Description: "Updated to kubernetes 1.8.4. Fixes a goroutine leak in the k8s api.",
							Kind:        versionbundle.KindChanged,
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "calico",
							Version: "2.6.2",
						},
						{
							Name:    "docker",
							Version: "1.12.6",
						},
						{
							Name:    "etcd",
							Version: "3.2.7",
						},
						{
							Name:    "kubedns",
							Version: "1.14.5",
						},
						{
							Name:    "kubernetes",
							Version: "1.8.4",
						},
						{
							Name:    "nginx-ingress-controller",
							Version: "0.9.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2017, time.December, 12, 10, 00, 0, 0, time.UTC),
					Version:      "1.0.0",
					WIP:          false,
				},
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "kubernetes",
							Description: "Enable encryption at rest",
							Kind:        versionbundle.KindChanged,
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "calico",
							Version: "2.6.2",
						},
						{
							Name:    "docker",
							Version: "1.12.6",
						},
						{
							Name:    "etcd",
							Version: "3.2.7",
						},
						{
							Name:    "kubedns",
							Version: "1.14.5",
						},
						{
							Name:    "kubernetes",
							Version: "1.8.4",
						},
						{
							Name:    "nginx-ingress-controller",
							Version: "0.9.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2017, time.December, 19, 10, 00, 0, 0, time.UTC),
					Version:      "1.1.0",
					WIP:          false,
				},
				{
					Changelogs: []versionbundle.Changelog{
						{

							Component:   "containerlinux",
							Description: "Updated containerlinux version to 1576.5.0.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "kubernetes",
							Description: "Fixed audit log.",
							Kind:        "fixed",
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "calico",
							Version: "2.6.2",
						},
						{
							Name:    "containerlinux",
							Version: "1576.5.0",
						},
						{
							Name:    "docker",
							Version: "17.09.0",
						},
						{
							Name:    "etcd",
							Version: "3.2.7",
						},
						{
							Name:    "kubedns",
							Version: "1.14.5",
						},
						{
							Name:    "kubernetes",
							Version: "1.8.4",
						},
						{
							Name:    "nginx-ingress-controller",
							Version: "0.9.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.February, 8, 6, 25, 0, 0, time.UTC),
					Version:      "1.2.0",
					WIP:          false,
				},
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "Kubernetes",
							Description: "Updated to Kubernetes 1.9.2.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "Kubernetes",
							Description: "Switched to vanilla (previously CoreOS) hyperkube image.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "Docker",
							Description: "Updated to 17.09.0-ce.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "Calico",
							Description: "Updated to 3.0.1.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "CoreDNS",
							Description: "Version 1.0.5 replaces kube-dns.",
							Kind:        versionbundle.KindAdded,
						},
						{
							Component:   "Nginx Ingress Controller",
							Description: "Updated to 0.10.2.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "cloudconfig",
							Description: "Add OIDC integration for Kubernetes api-server.",
							Kind:        versionbundle.KindAdded,
						},
						{
							Component:   "cloudconfig",
							Description: "Replace systemd units for Kubernetes components with self-hosted pods.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "containerlinux",
							Description: "Updated Container Linux version to 1576.5.0.",
							Kind:        versionbundle.KindChanged,
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "calico",
							Version: "3.0.1",
						},
						{
							Name:    "containerlinux",
							Version: "1576.5.0",
						},
						{
							Name:    "docker",
							Version: "17.09.0",
						},
						{
							Name:    "etcd",
							Version: "3.2.7",
						},
						{
							Name:    "coredns",
							Version: "1.0.5",
						},
						{
							Name:    "kubernetes",
							Version: "1.9.2",
						},
						{
							Name:    "nginx-ingress-controller",
							Version: "0.10.2",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.February, 15, 2, 27, 0, 0, time.UTC),
					Version:      "2.0.0",
					WIP:          false,
				},
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "kvm-node-controller",
							Description: "Updated KVM node controller with pod status bugfix.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "Calico",
							Description: "Updated to 3.0.2.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "kubelet",
							Description: "Tune kubelet flags for protecting key units (kubelet and container runtime) from workload overloads.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "etcd",
							Description: "Updated to 3.3.1.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "qemu",
							Description: "Fixed formula for calculating qemu memory overhead.",
							Kind:        versionbundle.KindFixed,
						},
						{
							Component:   "monitoring",
							Description: "Added configuration for monitoring endpoint IP addresses.",
							Kind:        versionbundle.KindAdded,
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "calico",
							Version: "3.0.2",
						},
						{
							Name:    "containerlinux",
							Version: "1576.5.0",
						},
						{
							Name:    "docker",
							Version: "17.09.0",
						},
						{
							Name:    "etcd",
							Version: "3.3.1",
						},
						{
							Name:    "coredns",
							Version: "1.0.5",
						},
						{
							Name:    "kubernetes",
							Version: "1.9.2",
						},
						{
							Name:    "nginx-ingress-controller",
							Version: "0.10.2",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.February, 20, 2, 57, 0, 0, time.UTC),
					Version:      "2.0.1",
					WIP:          false,
				},
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "kvm-node-controller",
							Description: "Updated KVM node controller with pod status bugfix.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "Calico",
							Description: "Updated to 3.0.2.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "kubelet",
							Description: "Tune kubelet flags for protecting key units (kubelet and container runtime) from workload overloads.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "etcd",
							Description: "Updated to 3.3.1.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "qemu",
							Description: "Fixed formula for calculating qemu memory overhead.",
							Kind:        versionbundle.KindFixed,
						},
						{
							Component:   "monitoring",
							Description: "Added configuration for monitoring endpoint IP addresses.",
							Kind:        versionbundle.KindAdded,
						},
						{
							Component:   "cloudconfig",
							Description: "Enable aggregation layer to be able to extend kubernetes API.",
							Kind:        versionbundle.KindChanged,
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "calico",
							Version: "3.0.2",
						},
						{
							Name:    "containerlinux",
							Version: "1576.5.0",
						},
						{
							Name:    "docker",
							Version: "17.09.0",
						},
						{
							Name:    "etcd",
							Version: "3.3.1",
						},
						{
							Name:    "coredns",
							Version: "1.0.6",
						},
						{
							Name:    "kubernetes",
							Version: "1.9.2",
						},
						{
							Name:    "nginx-ingress-controller",
							Version: "0.11.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.March, 1, 2, 57, 0, 0, time.UTC),
					Version:      "2.1.0",
					WIP:          false,
				},
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "containerlinux",
							Description: "Updated to version 1632.3.0.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "cloudconfig",
							Description: "Removed set-ownership-etcd-data-dir.service.",
							Kind:        versionbundle.KindChanged,
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "calico",
							Version: "3.0.2",
						},
						{
							Name:    "containerlinux",
							Version: "1632.3.0",
						},
						{
							Name:    "docker",
							Version: "17.09.0",
						},
						{
							Name:    "etcd",
							Version: "3.3.1",
						},
						{
							Name:    "coredns",
							Version: "1.0.6",
						},
						{
							Name:    "kubernetes",
							Version: "1.9.2",
						},
						{
							Name:    "nginx-ingress-controller",
							Version: "0.11.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.March, 7, 2, 57, 0, 0, time.UTC),
					Version:      "2.1.1",
					WIP:          false,
				},
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "cloudconfig",
							Description: "Kubernetes updated to version 1.9.5.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "cloudconfig",
							Description: "Nginx Ingress Controller updated to version 0.12.0",
							Kind:        versionbundle.KindChanged,
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "calico",
							Version: "3.0.2",
						},
						{
							Name:    "containerlinux",
							Version: "1632.3.0",
						},
						{
							Name:    "docker",
							Version: "17.09.0",
						},
						{
							Name:    "etcd",
							Version: "3.3.1",
						},
						{
							Name:    "coredns",
							Version: "1.0.6",
						},
						{
							Name:    "kubernetes",
							Version: "1.9.5",
						},
						{
							Name:    "nginx-ingress-controller",
							Version: "0.12.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.March, 13, 12, 30, 0, 0, time.UTC),
					Version:      "2.1.2",
					WIP:          false,
				},
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "cloudconfig",
							Description: "Improved encryption key injection.",
							Kind:        versionbundle.KindChanged,
						},
						{
							Component:   "kvmconfig",
							Description: "Updated ingress annotations for api and etcd.",
							Kind:        versionbundle.KindChanged,
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "calico",
							Version: "3.0.2",
						},
						{
							Name:    "containerlinux",
							Version: "1632.3.0",
						},
						{
							Name:    "docker",
							Version: "17.09.0",
						},
						{
							Name:    "etcd",
							Version: "3.3.1",
						},
						{
							Name:    "coredns",
							Version: "1.0.6",
						},
						{
							Name:    "kubernetes",
							Version: "1.9.5",
						},
						{
							Name:    "nginx-ingress-controller",
							Version: "0.12.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   false,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.March, 22, 15, 30, 0, 0, time.UTC),
					Version:      "2.1.3",
					WIP:          false,
				},
				{
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "components",
							Description: "Put your description here.",
							Kind:        versionbundle.KindChanged,
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "calico",
							Version: "3.0.2",
						},
						{
							Name:    "containerlinux",
							Version: "1632.3.0",
						},
						{
							Name:    "docker",
							Version: "17.09.0",
						},
						{
							Name:    "etcd",
							Version: "3.3.1",
						},
						{
							Name:    "coredns",
							Version: "1.0.6",
						},
						{
							Name:    "kubernetes",
							Version: "1.9.5",
						},
						{
							Name:    "nginx-ingress-controller",
							Version: "0.12.0",
						},
					},
					Dependencies: []versionbundle.Dependency{},
					Deprecated:   false,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.April, 04, 17, 32, 0, 0, time.UTC),
					Version:      "2.1.4",
					WIP:          true,
				},
			},
		}
		b, err := json.Marshal(cr)
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
		_, err = io.WriteString(w, string(b))
		if err != nil {
			t.Fatalf("expected %#v got %#v", nil, err)
		}
	}
}
