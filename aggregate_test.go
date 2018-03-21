package versionbundle

import (
	"reflect"
	"testing"
	"time"

	"github.com/giantswarm/micrologger/microloggertest"
	"github.com/kylelemons/godebug/pretty"
)

func Test_Aggregate(t *testing.T) {
	testCases := []struct {
		Bundles                []Bundle
		ExpectedGroupedBundles [][]Bundle
		ErrorMatcher           func(err error) bool
	}{
		// Test 0 ensures that nil input results in empty output.
		{
			Bundles:                nil,
			ExpectedGroupedBundles: nil,
			ErrorMatcher:           nil,
		},

		// Test 1 is the same as 0 but with an empty list of version bundles.
		{
			Bundles:                []Bundle{},
			ExpectedGroupedBundles: nil,
			ErrorMatcher:           nil,
		},

		// Test 2 ensures a single version bundle within the given list of version bundles
		// is within the aggregated state as it is.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
							Kind:        "changed",
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version requirements changed due to calico update.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kube-dns",
							Version: "1.0.0",
						},
					},
					Dependencies: []Dependency{
						{
							Name:    "kubernetes",
							Version: "<= 1.7.x",
						},
					},
					Deprecated: false,
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.1.0",
					WIP:        false,
				},
			},
			ExpectedGroupedBundles: [][]Bundle{
				{
					{
						Changelogs: []Changelog{
							{
								Component:   "calico",
								Description: "Calico version updated.",
								Kind:        "changed",
							},
							{
								Component:   "kubernetes",
								Description: "Kubernetes version requirements changed due to calico update.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							{
								Name:    "calico",
								Version: "1.1.0",
							},
							{
								Name:    "kube-dns",
								Version: "1.0.0",
							},
						},
						Dependencies: []Dependency{
							{
								Name:    "kubernetes",
								Version: "<= 1.7.x",
							},
						},
						Deprecated: false,
						Name:       "kubernetes-operator",
						Time:       time.Unix(10, 5),
						Version:    "0.1.0",
						WIP:        false,
					},
				},
			},
			ErrorMatcher: nil,
		},

		// Test 3 ensures depending version bundles within the given list of version
		// bundles are aggregated together within the aggregated state.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "etcd",
							Description: "Etcd version updated.",
							Kind:        "changed",
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "etcd",
							Version: "3.2.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.1",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(20, 15),
					Version:      "0.2.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
							Kind:        "changed",
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version requirements changed due to calico update.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kube-dns",
							Version: "1.0.0",
						},
					},
					Dependencies: []Dependency{
						{
							Name:    "kubernetes",
							Version: "<= 1.7.x",
						},
					},
					Deprecated: false,
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.1.0",
					WIP:        false,
				},
			},
			ExpectedGroupedBundles: [][]Bundle{
				{
					{
						Changelogs: []Changelog{
							{
								Component:   "etcd",
								Description: "Etcd version updated.",
								Kind:        "changed",
							},
							{
								Component:   "kubernetes",
								Description: "Kubernetes version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							{
								Name:    "etcd",
								Version: "3.2.0",
							},
							{
								Name:    "kubernetes",
								Version: "1.7.1",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cloud-config-operator",
						Time:         time.Unix(20, 15),
						Version:      "0.2.0",
						WIP:          false,
					},
					{
						Changelogs: []Changelog{
							{
								Component:   "calico",
								Description: "Calico version updated.",
								Kind:        "changed",
							},
							{
								Component:   "kubernetes",
								Description: "Kubernetes version requirements changed due to calico update.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							{
								Name:    "calico",
								Version: "1.1.0",
							},
							{
								Name:    "kube-dns",
								Version: "1.0.0",
							},
						},
						Dependencies: []Dependency{
							{
								Name:    "kubernetes",
								Version: "<= 1.7.x",
							},
						},
						Deprecated: false,
						Name:       "kubernetes-operator",
						Time:       time.Unix(10, 5),
						Version:    "0.1.0",
						WIP:        false,
					},
				},
			},
			ErrorMatcher: nil,
		},

		// Test 4 ensures depending version bundles are not aggregated together in
		// case their dependency definitions do not meet the defined constraints.
		// Thus the aggregated state should be empty because there is no proper
		// release available.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "etcd",
							Description: "Etcd version updated.",
							Kind:        "changed",
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "etcd",
							Version: "3.2.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.1",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(20, 15),
					Version:      "0.2.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
							Kind:        "changed",
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version requirements changed due to calico update.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kube-dns",
							Version: "1.0.0",
						},
					},
					Dependencies: []Dependency{
						{
							Name:    "kubernetes",
							Version: "<= 1.7.0",
						},
					},
					Deprecated: false,
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.1.0",
					WIP:        false,
				},
			},
			ExpectedGroupedBundles: nil,
			ErrorMatcher:           nil,
		},

		// Test 5 ensures when having an operator's version bundles [a1,a2] and
		// having another operator's version bundles [b1], there should be
		// two aggregated releases [[a1,b1],[a2,b1]].
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "etcd",
							Description: "Etcd version updated.",
							Kind:        "changed",
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "etcd",
							Version: "3.2.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.1",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(20, 15),
					Version:      "0.2.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "kubernetes",
							Description: "Kubernetes version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "etcd",
							Version: "3.2.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.8.1",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(30, 20),
					Version:      "0.3.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
							Kind:        "changed",
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version requirements changed due to calico update.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kube-dns",
							Version: "1.0.0",
						},
					},
					Dependencies: []Dependency{
						{
							Name:    "kubernetes",
							Version: "<= 1.8.x",
						},
					},
					Deprecated: false,
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.1.0",
					WIP:        false,
				},
			},
			ExpectedGroupedBundles: [][]Bundle{
				{
					{
						Changelogs: []Changelog{
							{
								Component:   "etcd",
								Description: "Etcd version updated.",
								Kind:        "changed",
							},
							{
								Component:   "kubernetes",
								Description: "Kubernetes version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							{
								Name:    "etcd",
								Version: "3.2.0",
							},
							{
								Name:    "kubernetes",
								Version: "1.7.1",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cloud-config-operator",
						Time:         time.Unix(20, 15),
						Version:      "0.2.0",
						WIP:          false,
					},
					{
						Changelogs: []Changelog{
							{
								Component:   "calico",
								Description: "Calico version updated.",
								Kind:        "changed",
							},
							{
								Component:   "kubernetes",
								Description: "Kubernetes version requirements changed due to calico update.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							{
								Name:    "calico",
								Version: "1.1.0",
							},
							{
								Name:    "kube-dns",
								Version: "1.0.0",
							},
						},
						Dependencies: []Dependency{
							{
								Name:    "kubernetes",
								Version: "<= 1.8.x",
							},
						},
						Deprecated: false,
						Name:       "kubernetes-operator",
						Time:       time.Unix(10, 5),
						Version:    "0.1.0",
						WIP:        false,
					},
				},
				{
					{
						Changelogs: []Changelog{
							{
								Component:   "kubernetes",
								Description: "Kubernetes version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							{
								Name:    "etcd",
								Version: "3.2.0",
							},
							{
								Name:    "kubernetes",
								Version: "1.8.1",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cloud-config-operator",
						Time:         time.Unix(30, 20),
						Version:      "0.3.0",
						WIP:          false,
					},
					{
						Changelogs: []Changelog{
							{
								Component:   "calico",
								Description: "Calico version updated.",
								Kind:        "changed",
							},
							{
								Component:   "kubernetes",
								Description: "Kubernetes version requirements changed due to calico update.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							{
								Name:    "calico",
								Version: "1.1.0",
							},
							{
								Name:    "kube-dns",
								Version: "1.0.0",
							},
						},
						Dependencies: []Dependency{
							{
								Name:    "kubernetes",
								Version: "<= 1.8.x",
							},
						},
						Deprecated: false,
						Name:       "kubernetes-operator",
						Time:       time.Unix(10, 5),
						Version:    "0.1.0",
						WIP:        false,
					},
				},
			},
			ErrorMatcher: nil,
		},

		// Test 6 ensures when having an operator's version bundles [a1,a2] and
		// having another operator's version bundles [b1], there should be one
		// aggregated release [[a2,b1]].
		//
		// NOTE a1 requires a dependency which cannot be fulfilled. This is why
		// there is only one possible release.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "etcd",
							Description: "Etcd version updated.",
							Kind:        "changed",
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "etcd",
							Version: "3.2.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.1",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(20, 15),
					Version:      "0.2.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "kubernetes",
							Description: "Kubernetes version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "etcd",
							Version: "3.2.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.8.1",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cloud-config-operator",
					Time:         time.Unix(30, 20),
					Version:      "0.3.0",
					WIP:          false,
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "calico",
							Description: "Calico version updated.",
							Kind:        "changed",
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version requirements changed due to calico update.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kube-dns",
							Version: "1.0.0",
						},
					},
					Dependencies: []Dependency{
						{
							Name:    "kubernetes",
							Version: "== 1.8.1",
						},
					},
					Deprecated: false,
					Name:       "kubernetes-operator",
					Time:       time.Unix(10, 5),
					Version:    "0.1.0",
					WIP:        false,
				},
			},
			ExpectedGroupedBundles: [][]Bundle{
				{
					{
						Changelogs: []Changelog{
							{
								Component:   "kubernetes",
								Description: "Kubernetes version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							{
								Name:    "etcd",
								Version: "3.2.0",
							},
							{
								Name:    "kubernetes",
								Version: "1.8.1",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cloud-config-operator",
						Time:         time.Unix(30, 20),
						Version:      "0.3.0",
						WIP:          false,
					},
					{
						Changelogs: []Changelog{
							{
								Component:   "calico",
								Description: "Calico version updated.",
								Kind:        "changed",
							},
							{
								Component:   "kubernetes",
								Description: "Kubernetes version requirements changed due to calico update.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							{
								Name:    "calico",
								Version: "1.1.0",
							},
							{
								Name:    "kube-dns",
								Version: "1.0.0",
							},
						},
						Dependencies: []Dependency{
							{
								Name:    "kubernetes",
								Version: "== 1.8.1",
							},
						},
						Deprecated: false,
						Name:       "kubernetes-operator",
						Time:       time.Unix(10, 5),
						Version:    "0.1.0",
						WIP:        false,
					},
				},
			},
			ErrorMatcher: nil,
		},

		// Test 7 replicates a situation in geckon where incorrect grouped bundles where returned
		{
			Bundles: []Bundle{
				Bundle{
					Changelogs: []Changelog{
						Changelog{
							Component:   "vault",
							Description: "Vault version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						Component{
							Name:    "vault",
							Version: "0.7.3",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "cert-operator",
					Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
					Version:      "0.1.0",
					WIP:          false,
				},
				Bundle{
					Changelogs: []Changelog{
						Changelog{
							Component:   "flannel",
							Description: "Flannel version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						Component{
							Name:    "flannel",
							Version: "0.9.0",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "flannel-operator",
					Time:         time.Date(2017, time.October, 27, 16, 21, 0, 0, time.UTC),
					Version:      "0.1.0",
					WIP:          false,
				},
				Bundle{
					Changelogs: []Changelog{
						Changelog{
							Component:   "flannel",
							Description: "Flannel version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						Component{
							Name:    "flannel",
							Version: "0.9.0",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "flannel-operator",
					Time:         time.Date(2018, time.March, 16, 9, 15, 0, 0, time.UTC),
					Version:      "0.2.0",
					WIP:          true,
				},
				Bundle{
					Changelogs: []Changelog{
						Changelog{
							Component:   "calico",
							Description: "Calico version updated.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "docker",
							Description: "Docker version updated.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "etcd",
							Description: "Etcd version updated.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "kubedns",
							Description: "KubeDNS version updated.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "kubernetes",
							Description: "Kubernetes version updated.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "nginx-ingress-controller",
							Description: "Nginx-ingress-controller version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						Component{
							Name:    "calico",
							Version: "2.6.2",
						},
						Component{
							Name:    "docker",
							Version: "1.12.6",
						},
						Component{
							Name:    "etcd",
							Version: "3.2.7",
						},
						Component{
							Name:    "kubedns",
							Version: "1.14.5",
						},
						Component{
							Name:    "kubernetes",
							Version: "1.8.1",
						},
						Component{
							Name:    "nginx-ingress-controller",
							Version: "0.9.0",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2017, time.October, 26, 16, 38, 0, 0, time.UTC),
					Version:      "0.1.0",
					WIP:          false,
				},
				Bundle{
					Changelogs: []Changelog{
						Changelog{
							Component:   "kubernetes",
							Description: "Updated to kubernetes 1.8.4. Fixes a goroutine leak in the k8s api.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						Component{
							Name:    "calico",
							Version: "2.6.2",
						},
						Component{
							Name:    "docker",
							Version: "1.12.6",
						},
						Component{
							Name:    "etcd",
							Version: "3.2.7",
						},
						Component{
							Name:    "kubedns",
							Version: "1.14.5",
						},
						Component{
							Name:    "kubernetes",
							Version: "1.8.4",
						},
						Component{
							Name:    "nginx-ingress-controller",
							Version: "0.9.0",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2017, time.December, 12, 10, 0, 0, 0, time.UTC),
					Version:      "1.0.0",
					WIP:          false,
				},
				Bundle{
					Changelogs: []Changelog{
						Changelog{
							Component:   "kubernetes",
							Description: "Enable encryption at rest",
							Kind:        "changed",
						},
					},
					Components: []Component{
						Component{
							Name:    "calico",
							Version: "2.6.2",
						},
						Component{
							Name:    "docker",
							Version: "1.12.6",
						},
						Component{
							Name:    "etcd",
							Version: "3.2.7",
						},
						Component{
							Name:    "kubedns",
							Version: "1.14.5",
						},
						Component{
							Name:    "kubernetes",
							Version: "1.8.4",
						},
						Component{
							Name:    "nginx-ingress-controller",
							Version: "0.9.0",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2017, time.December, 19, 10, 0, 0, 0, time.UTC),
					Version:      "1.1.0",
					WIP:          false,
				},
				Bundle{
					Changelogs: []Changelog{
						Changelog{
							Component:   "containerlinux",
							Description: "Updated containerlinux version to 1576.5.0.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "kubernetes",
							Description: "Fixed audit log.",
							Kind:        "fixed",
						},
					},
					Components: []Component{
						Component{
							Name:    "calico",
							Version: "2.6.2",
						},
						Component{
							Name:    "containerlinux",
							Version: "1576.5.0",
						},
						Component{
							Name:    "docker",
							Version: "17.09.0",
						},
						Component{
							Name:    "etcd",
							Version: "3.2.7",
						},
						Component{
							Name:    "kubedns",
							Version: "1.14.5",
						},
						Component{
							Name:    "kubernetes",
							Version: "1.8.4",
						},
						Component{
							Name:    "nginx-ingress-controller",
							Version: "0.9.0",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.February, 8, 6, 25, 0, 0, time.UTC),
					Version:      "1.2.0",
					WIP:          false,
				},
				Bundle{
					Changelogs: []Changelog{
						Changelog{
							Component:   "Kubernetes",
							Description: "Updated to Kubernetes 1.9.2.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "Kubernetes",
							Description: "Switched to vanilla (previously CoreOS) hyperkube image.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "Docker",
							Description: "Updated to 17.09.0-ce.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "Calico",
							Description: "Updated to 3.0.1.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "CoreDNS",
							Description: "Version 1.0.5 replaces kube-dns.",
							Kind:        "added",
						},
						Changelog{
							Component:   "Nginx Ingress Controller",
							Description: "Updated to 0.10.2.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "cloudconfig",
							Description: "Add OIDC integration for Kubernetes api-server.",
							Kind:        "added",
						},
						Changelog{
							Component:   "cloudconfig",
							Description: "Replace systemd units for Kubernetes components with self-hosted pods.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "containerlinux",
							Description: "Updated Container Linux version to 1576.5.0.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						Component{
							Name:    "calico",
							Version: "3.0.1",
						},
						Component{
							Name:    "containerlinux",
							Version: "1576.5.0",
						},
						Component{
							Name:    "docker",
							Version: "17.09.0",
						},
						Component{
							Name:    "etcd",
							Version: "3.2.7",
						},
						Component{
							Name:    "coredns",
							Version: "1.0.5",
						},
						Component{
							Name:    "kubernetes",
							Version: "1.9.2",
						},
						Component{
							Name:    "nginx-ingress-controller",
							Version: "0.10.2",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.February, 15, 2, 27, 0, 0, time.UTC),
					Version:      "2.0.0",
					WIP:          false,
				},
				Bundle{
					Changelogs: []Changelog{
						Changelog{
							Component:   "kvm-node-controller",
							Description: "Updated KVM node controller with pod status bugfix.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "Calico",
							Description: "Updated to 3.0.2.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "kubelet",
							Description: "Tune kubelet flags for protecting key units (kubelet and container runtime) from workload overloads.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "etcd",
							Description: "Updated to 3.3.1.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "qemu",
							Description: "Fixed formula for calculating qemu memory overhead.",
							Kind:        "fixed",
						},
						Changelog{
							Component:   "monitoring",
							Description: "Added configuration for monitoring endpoint IP addresses.",
							Kind:        "added",
						},
					},
					Components: []Component{
						Component{
							Name:    "calico",
							Version: "3.0.2",
						},
						Component{
							Name:    "containerlinux",
							Version: "1576.5.0",
						},
						Component{
							Name:    "docker",
							Version: "17.09.0",
						},
						Component{
							Name:    "etcd",
							Version: "3.3.1",
						},
						Component{
							Name:    "coredns",
							Version: "1.0.5",
						},
						Component{
							Name:    "kubernetes",
							Version: "1.9.2",
						},
						Component{
							Name:    "nginx-ingress-controller",
							Version: "0.10.2",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.February, 20, 2, 57, 0, 0, time.UTC),
					Version:      "2.0.1",
					WIP:          false,
				},
				Bundle{
					Changelogs: []Changelog{
						Changelog{
							Component:   "kvm-node-controller",
							Description: "Updated KVM node controller with pod status bugfix.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "Calico",
							Description: "Updated to 3.0.2.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "kubelet",
							Description: "Tune kubelet flags for protecting key units (kubelet and container runtime) from workload overloads.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "etcd",
							Description: "Updated to 3.3.1.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "qemu",
							Description: "Fixed formula for calculating qemu memory overhead.",
							Kind:        "fixed",
						},
						Changelog{
							Component:   "monitoring",
							Description: "Added configuration for monitoring endpoint IP addresses.",
							Kind:        "added",
						},
						Changelog{
							Component:   "cloudconfig",
							Description: "Enable aggregation layer to be able to extend kubernetes API.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						Component{
							Name:    "calico",
							Version: "3.0.2",
						},
						Component{
							Name:    "containerlinux",
							Version: "1576.5.0",
						},
						Component{
							Name:    "docker",
							Version: "17.09.0",
						},
						Component{
							Name:    "etcd",
							Version: "3.3.1",
						},
						Component{
							Name:    "coredns",
							Version: "1.0.6",
						},
						Component{
							Name:    "kubernetes",
							Version: "1.9.2",
						},
						Component{
							Name:    "nginx-ingress-controller",
							Version: "0.11.0",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   true,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.March, 1, 2, 57, 0, 0, time.UTC),
					Version:      "2.1.0",
					WIP:          false,
				},
				Bundle{
					Changelogs: []Changelog{
						Changelog{
							Component:   "containerlinux",
							Description: "Updated to version 1632.3.0.",
							Kind:        "changed",
						},
						Changelog{
							Component:   "cloudconfig",
							Description: "Removed set-ownership-etcd-data-dir.service.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						Component{
							Name:    "calico",
							Version: "3.0.2",
						},
						Component{
							Name:    "containerlinux",
							Version: "1632.3.0",
						},
						Component{
							Name:    "docker",
							Version: "17.09.0",
						},
						Component{
							Name:    "etcd",
							Version: "3.3.1",
						},
						Component{
							Name:    "coredns",
							Version: "1.0.6",
						},
						Component{
							Name:    "kubernetes",
							Version: "1.9.2",
						},
						Component{
							Name:    "nginx-ingress-controller",
							Version: "0.11.0",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.March, 7, 2, 57, 0, 0, time.UTC),
					Version:      "2.1.1",
					WIP:          false,
				},
				Bundle{
					Changelogs: []Changelog{
						Changelog{
							Component:   "component",
							Description: "Put your description here",
							Kind:        "changed",
						},
					},
					Components: []Component{
						Component{
							Name:    "calico",
							Version: "3.0.2",
						},
						Component{
							Name:    "containerlinux",
							Version: "1632.3.0",
						},
						Component{
							Name:    "docker",
							Version: "17.09.0",
						},
						Component{
							Name:    "etcd",
							Version: "3.3.1",
						},
						Component{
							Name:    "coredns",
							Version: "1.0.6",
						},
						Component{
							Name:    "kubernetes",
							Version: "1.9.2",
						},
						Component{
							Name:    "nginx-ingress-controller",
							Version: "0.11.0",
						},
					},
					Dependencies: []Dependency{},
					Deprecated:   false,
					Name:         "kvm-operator",
					Time:         time.Date(2018, time.March, 13, 12, 30, 0, 0, time.UTC),
					Version:      "2.1.2",
					WIP:          true},
			},
			ExpectedGroupedBundles: [][]Bundle{
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2017, time.October, 27, 16, 21, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "calico",
								Description: "Calico version updated.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "docker",
								Description: "Docker version updated.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "etcd",
								Description: "Etcd version updated.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "kubedns",
								Description: "KubeDNS version updated.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "kubernetes",
								Description: "Kubernetes version updated.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "nginx-ingress-controller",
								Description: "Nginx-ingress-controller version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "2.6.2",
							},
							Component{
								Name:    "docker",
								Version: "1.12.6",
							},
							Component{
								Name:    "etcd",
								Version: "3.2.7",
							},
							Component{
								Name:    "kubedns",
								Version: "1.14.5",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.8.1",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2017, time.October, 26, 16, 38, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2017, time.October, 27, 16, 21, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "kubernetes",
								Description: "Updated to kubernetes 1.8.4. Fixes a goroutine leak in the k8s api.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "2.6.2",
							},
							Component{
								Name:    "docker",
								Version: "1.12.6",
							},
							Component{
								Name:    "etcd",
								Version: "3.2.7",
							},
							Component{
								Name:    "kubedns",
								Version: "1.14.5",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.8.4",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2017, time.December, 12, 10, 0, 0, 0, time.UTC),
						Version:      "1.0.0",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2017, 10, 27, 16, 21, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "kubernetes",
								Description: "Enable encryption at rest",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "2.6.2",
							},
							Component{
								Name:    "docker",
								Version: "1.12.6",
							},
							Component{
								Name:    "etcd",
								Version: "3.2.7",
							},
							Component{
								Name:    "kubedns",
								Version: "1.14.5",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.8.4",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2017, time.December, 19, 10, 0, 0, 0, time.UTC),
						Version:      "1.1.0",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2017, time.October, 27, 16, 21, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "containerlinux",
								Description: "Updated containerlinux version to 1576.5.0.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "kubernetes",
								Description: "Fixed audit log.",
								Kind:        "fixed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "2.6.2",
							},
							Component{
								Name:    "containerlinux",
								Version: "1576.5.0",
							},
							Component{
								Name:    "docker",
								Version: "17.09.0",
							},
							Component{
								Name:    "etcd",
								Version: "3.2.7",
							},
							Component{
								Name:    "kubedns",
								Version: "1.14.5",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.8.4",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2018, time.February, 8, 6, 25, 0, 0, time.UTC),
						Version:      "1.2.0",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2017, time.October, 27, 16, 21, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "Kubernetes",
								Description: "Updated to Kubernetes 1.9.2.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "Kubernetes",
								Description: "Switched to vanilla (previously CoreOS) hyperkube image.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "Docker",
								Description: "Updated to 17.09.0-ce.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "Calico",
								Description: "Updated to 3.0.1.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "CoreDNS",
								Description: "Version 1.0.5 replaces kube-dns.",
								Kind:        "added",
							},
							Changelog{
								Component:   "Nginx Ingress Controller",
								Description: "Updated to 0.10.2.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "cloudconfig",
								Description: "Add OIDC integration for Kubernetes api-server.",
								Kind:        "added",
							},
							Changelog{
								Component:   "cloudconfig",
								Description: "Replace systemd units for Kubernetes components with self-hosted pods.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "containerlinux",
								Description: "Updated Container Linux version to 1576.5.0.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "3.0.1",
							},
							Component{
								Name:    "containerlinux",
								Version: "1576.5.0",
							},
							Component{
								Name:    "docker",
								Version: "17.09.0",
							},
							Component{
								Name:    "etcd",
								Version: "3.2.7",
							},
							Component{
								Name:    "coredns",
								Version: "1.0.5",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.9.2",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.10.2",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2018, time.February, 15, 2, 27, 0, 0, time.UTC),
						Version:      "2.0.0",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2017, time.October, 27, 16, 21, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "kvm-node-controller",
								Description: "Updated KVM node controller with pod status bugfix.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "Calico",
								Description: "Updated to 3.0.2.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "kubelet",
								Description: "Tune kubelet flags for protecting key units (kubelet and container runtime) from workload overloads.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "etcd",
								Description: "Updated to 3.3.1.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "qemu",
								Description: "Fixed formula for calculating qemu memory overhead.",
								Kind:        "fixed",
							},
							Changelog{
								Component:   "monitoring",
								Description: "Added configuration for monitoring endpoint IP addresses.",
								Kind:        "added",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "3.0.2",
							},
							Component{
								Name:    "containerlinux",
								Version: "1576.5.0",
							},
							Component{
								Name:    "docker",
								Version: "17.09.0",
							},
							Component{
								Name:    "etcd",
								Version: "3.3.1",
							},
							Component{
								Name:    "coredns",
								Version: "1.0.5",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.9.2",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.10.2",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2018, time.February, 20, 2, 57, 0, 0, time.UTC),
						Version:      "2.0.1",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2017, time.October, 27, 16, 21, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "kvm-node-controller",
								Description: "Updated KVM node controller with pod status bugfix.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "Calico",
								Description: "Updated to 3.0.2.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "kubelet",
								Description: "Tune kubelet flags for protecting key units (kubelet and container runtime) from workload overloads.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "etcd",
								Description: "Updated to 3.3.1.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "qemu",
								Description: "Fixed formula for calculating qemu memory overhead.",
								Kind:        "fixed",
							},
							Changelog{
								Component:   "monitoring",
								Description: "Added configuration for monitoring endpoint IP addresses.",
								Kind:        "added",
							},
							Changelog{
								Component:   "cloudconfig",
								Description: "Enable aggregation layer to be able to extend kubernetes API.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "3.0.2",
							},
							Component{
								Name:    "containerlinux",
								Version: "1576.5.0",
							},
							Component{
								Name:    "docker",
								Version: "17.09.0",
							},
							Component{
								Name:    "etcd",
								Version: "3.3.1",
							},
							Component{
								Name:    "coredns",
								Version: "1.0.6",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.9.2",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.11.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2018, time.March, 1, 2, 57, 0, 0, time.UTC),
						Version:      "2.1.0",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2017, time.October, 27, 16, 21, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "containerlinux",
								Description: "Updated to version 1632.3.0.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "cloudconfig",
								Description: "Removed set-ownership-etcd-data-dir.service.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "3.0.2",
							},
							Component{
								Name:    "containerlinux",
								Version: "1632.3.0",
							},
							Component{
								Name:    "docker",
								Version: "17.09.0",
							},
							Component{
								Name:    "etcd",
								Version: "3.3.1",
							},
							Component{
								Name:    "coredns",
								Version: "1.0.6",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.9.2",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.11.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "kvm-operator",
						Time:         time.Date(2018, time.March, 7, 2, 57, 0, 0, time.UTC),
						Version:      "2.1.1",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2017, time.October, 27, 16, 21, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "component",
								Description: "Put your description here",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "3.0.2",
							},
							Component{
								Name:    "containerlinux",
								Version: "1632.3.0",
							},
							Component{
								Name:    "docker",
								Version: "17.09.0",
							},
							Component{
								Name:    "etcd",
								Version: "3.3.1",
							},
							Component{
								Name:    "coredns",
								Version: "1.0.6",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.9.2",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.11.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "kvm-operator",
						Time:         time.Date(2018, time.March, 13, 12, 30, 0, 0, time.UTC),
						Version:      "2.1.2",
						WIP:          true,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2018, time.March, 16, 9, 15, 0, 0, time.UTC),
						Version:      "0.2.0",
						WIP:          true,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "calico",
								Description: "Calico version updated.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "docker",
								Description: "Docker version updated.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "etcd",
								Description: "Etcd version updated.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "kubedns",
								Description: "KubeDNS version updated.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "kubernetes",
								Description: "Kubernetes version updated.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "nginx-ingress-controller",
								Description: "Nginx-ingress-controller version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "2.6.2",
							},
							Component{
								Name:    "docker",
								Version: "1.12.6",
							},
							Component{
								Name:    "etcd",
								Version: "3.2.7",
							},
							Component{
								Name:    "kubedns",
								Version: "1.14.5",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.8.1",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2017, time.October, 26, 16, 38, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2018, time.March, 16, 9, 15, 0, 0, time.UTC),
						Version:      "0.2.0",
						WIP:          true,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "kubernetes",
								Description: "Updated to kubernetes 1.8.4. Fixes a goroutine leak in the k8s api.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "2.6.2",
							},
							Component{
								Name:    "docker",
								Version: "1.12.6",
							},
							Component{
								Name:    "etcd",
								Version: "3.2.7",
							},
							Component{
								Name:    "kubedns",
								Version: "1.14.5",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.8.4",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2017, time.December, 12, 10, 0, 0, 0, time.UTC),
						Version:      "1.0.0",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2018, time.March, 16, 9, 15, 0, 0, time.UTC),
						Version:      "0.2.0",
						WIP:          true,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "kubernetes",
								Description: "Enable encryption at rest",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "2.6.2",
							},
							Component{
								Name:    "docker",
								Version: "1.12.6",
							},
							Component{
								Name:    "etcd",
								Version: "3.2.7",
							},
							Component{
								Name:    "kubedns",
								Version: "1.14.5",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.8.4",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2017, time.December, 19, 10, 0, 0, 0, time.UTC),
						Version:      "1.1.0",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2018, time.March, 16, 9, 15, 0, 0, time.UTC),
						Version:      "0.2.0",
						WIP:          true,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "containerlinux",
								Description: "Updated containerlinux version to 1576.5.0.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "kubernetes",
								Description: "Fixed audit log.",
								Kind:        "fixed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "2.6.2",
							},
							Component{
								Name:    "containerlinux",
								Version: "1576.5.0",
							},
							Component{
								Name:    "docker",
								Version: "17.09.0",
							},
							Component{
								Name:    "etcd",
								Version: "3.2.7",
							},
							Component{
								Name:    "kubedns",
								Version: "1.14.5",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.8.4",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2018, time.February, 8, 6, 25, 0, 0, time.UTC),
						Version:      "1.2.0",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2018, time.March, 16, 9, 15, 0, 0, time.UTC),
						Version:      "0.2.0",
						WIP:          true,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "Kubernetes",
								Description: "Updated to Kubernetes 1.9.2.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "Kubernetes",
								Description: "Switched to vanilla (previously CoreOS) hyperkube image.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "Docker",
								Description: "Updated to 17.09.0-ce.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "Calico",
								Description: "Updated to 3.0.1.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "CoreDNS",
								Description: "Version 1.0.5 replaces kube-dns.",
								Kind:        "added",
							},
							Changelog{
								Component:   "Nginx Ingress Controller",
								Description: "Updated to 0.10.2.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "cloudconfig",
								Description: "Add OIDC integration for Kubernetes api-server.",
								Kind:        "added",
							},
							Changelog{
								Component:   "cloudconfig",
								Description: "Replace systemd units for Kubernetes components with self-hosted pods.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "containerlinux",
								Description: "Updated Container Linux version to 1576.5.0.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "3.0.1",
							},
							Component{
								Name:    "containerlinux",
								Version: "1576.5.0",
							},
							Component{
								Name:    "docker",
								Version: "17.09.0",
							},
							Component{
								Name:    "etcd",
								Version: "3.2.7",
							},
							Component{
								Name:    "coredns",
								Version: "1.0.5",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.9.2",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.10.2",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2018, time.February, 15, 2, 27, 0, 0, time.UTC),
						Version:      "2.0.0",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2018, time.March, 16, 9, 15, 0, 0, time.UTC),
						Version:      "0.2.0",
						WIP:          true,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "kvm-node-controller",
								Description: "Updated KVM node controller with pod status bugfix.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "Calico",
								Description: "Updated to 3.0.2.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "kubelet",
								Description: "Tune kubelet flags for protecting key units (kubelet and container runtime) from workload overloads.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "etcd",
								Description: "Updated to 3.3.1.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "qemu",
								Description: "Fixed formula for calculating qemu memory overhead.",
								Kind:        "fixed",
							},
							Changelog{
								Component:   "monitoring",
								Description: "Added configuration for monitoring endpoint IP addresses.",
								Kind:        "added",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "3.0.2",
							},
							Component{
								Name:    "containerlinux",
								Version: "1576.5.0",
							},
							Component{
								Name:    "docker",
								Version: "17.09.0",
							},
							Component{
								Name:    "etcd",
								Version: "3.3.1",
							},
							Component{
								Name:    "coredns",
								Version: "1.0.5",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.9.2",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.10.2",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2018, time.February, 20, 2, 57, 0, 0, time.UTC),
						Version:      "2.0.1",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2018, time.March, 16, 9, 15, 0, 0, time.UTC),
						Version:      "0.2.0",
						WIP:          true,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "kvm-node-controller",
								Description: "Updated KVM node controller with pod status bugfix.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "Calico",
								Description: "Updated to 3.0.2.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "kubelet",
								Description: "Tune kubelet flags for protecting key units (kubelet and container runtime) from workload overloads.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "etcd",
								Description: "Updated to 3.3.1.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "qemu",
								Description: "Fixed formula for calculating qemu memory overhead.",
								Kind:        "fixed",
							},
							Changelog{
								Component:   "monitoring",
								Description: "Added configuration for monitoring endpoint IP addresses.",
								Kind:        "added",
							},
							Changelog{
								Component:   "cloudconfig",
								Description: "Enable aggregation layer to be able to extend kubernetes API.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "3.0.2",
							},
							Component{
								Name:    "containerlinux",
								Version: "1576.5.0",
							},
							Component{
								Name:    "docker",
								Version: "17.09.0",
							},
							Component{
								Name:    "etcd",
								Version: "3.3.1",
							},
							Component{
								Name:    "coredns",
								Version: "1.0.6",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.9.2",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.11.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   true,
						Name:         "kvm-operator",
						Time:         time.Date(2018, time.March, 1, 2, 57, 0, 0, time.UTC),
						Version:      "2.1.0",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2018, time.March, 16, 9, 15, 0, 0, time.UTC),
						Version:      "0.2.0",
						WIP:          true,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "containerlinux",
								Description: "Updated to version 1632.3.0.",
								Kind:        "changed",
							},
							Changelog{
								Component:   "cloudconfig",
								Description: "Removed set-ownership-etcd-data-dir.service.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "3.0.2",
							},
							Component{
								Name:    "containerlinux",
								Version: "1632.3.0",
							},
							Component{
								Name:    "docker",
								Version: "17.09.0",
							},
							Component{
								Name:    "etcd",
								Version: "3.3.1",
							},
							Component{
								Name:    "coredns",
								Version: "1.0.6",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.9.2",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.11.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "kvm-operator",
						Time:         time.Date(2018, time.March, 7, 2, 57, 0, 0, time.UTC),
						Version:      "2.1.1",
						WIP:          false,
					},
				},
				[]Bundle{
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "vault",
								Description: "Vault version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "vault",
								Version: "0.7.3",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "cert-operator",
						Time:         time.Date(2017, time.October, 26, 16, 53, 0, 0, time.UTC),
						Version:      "0.1.0",
						WIP:          false,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "flannel",
								Description: "Flannel version updated.",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "flannel",
								Version: "0.9.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "flannel-operator",
						Time:         time.Date(2018, time.March, 16, 9, 15, 0, 0, time.UTC),
						Version:      "0.2.0",
						WIP:          true,
					},
					Bundle{
						Changelogs: []Changelog{
							Changelog{
								Component:   "component",
								Description: "Put your description here",
								Kind:        "changed",
							},
						},
						Components: []Component{
							Component{
								Name:    "calico",
								Version: "3.0.2",
							},
							Component{
								Name:    "containerlinux",
								Version: "1632.3.0",
							},
							Component{
								Name:    "docker",
								Version: "17.09.0",
							},
							Component{
								Name:    "etcd",
								Version: "3.3.1",
							},
							Component{
								Name:    "coredns",
								Version: "1.0.6",
							},
							Component{
								Name:    "kubernetes",
								Version: "1.9.2",
							},
							Component{
								Name:    "nginx-ingress-controller",
								Version: "0.11.0",
							},
						},
						Dependencies: []Dependency{},
						Deprecated:   false,
						Name:         "kvm-operator",
						Time:         time.Date(2018, time.March, 13, 12, 30, 0, 0, time.UTC),
						Version:      "2.1.2",
						WIP:          true,
					},
				},
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		var err error

		var a *Aggregator
		{
			c := AggregatorConfig{
				Logger: microloggertest.New(),
			}

			a, err = NewAggregator(c)
			if err != nil {
				t.Fatalf("test %d expected %#v got %#v", i, nil, err)
			}
		}

		groupedBundles, err := a.Aggregate(tc.Bundles)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else {
			if !reflect.DeepEqual(groupedBundles, tc.ExpectedGroupedBundles) {
				diff := pretty.Compare(tc.ExpectedGroupedBundles, groupedBundles)
				t.Fatalf("test %d expected %#v got %#v (\n%s \n)", i, tc.ExpectedGroupedBundles, groupedBundles, diff)
			}
		}
	}
}
