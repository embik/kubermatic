/*
Copyright 2020 The Kubermatic Kubernetes Platform contributors.

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

package metricsserver

import (
	"k8c.io/kubermatic/v2/pkg/resources"
	"k8c.io/reconciler/pkg/reconciling"

	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterRoleBindingResourceReaderReconciler returns the ClusterRoleBinding required for the metrics server to read all required resources.
func ClusterRoleBindingResourceReaderReconciler() reconciling.NamedClusterRoleBindingReconcilerFactory {
	return func() (string, reconciling.ClusterRoleBindingReconciler) {
		return resources.MetricsServerResourceReaderClusterRoleBindingName, func(crb *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
			crb.Labels = resources.BaseAppLabels(Name, nil)

			crb.RoleRef = rbacv1.RoleRef{
				Name:     resources.MetricsServerClusterRoleName,
				Kind:     "ClusterRole",
				APIGroup: rbacv1.GroupName,
			}

			// metrics server running in the user cluster - ServiceAccount
			crb.Subjects = []rbacv1.Subject{
				{
					Kind:      rbacv1.ServiceAccountKind,
					Name:      resources.MetricsServerServiceAccountName,
					Namespace: metav1.NamespaceSystem,
				},
			}

			return crb, nil
		}
	}
}

// ClusterRoleBindingAuthDelegatorReconciler returns the ClusterRoleBinding required for the metrics server to create token review requests.
func ClusterRoleBindingAuthDelegatorReconciler() reconciling.NamedClusterRoleBindingReconcilerFactory {
	return func() (string, reconciling.ClusterRoleBindingReconciler) {
		// metrics server running in the user cluster
		return "metrics-server:system:auth-delegator", func(crb *rbacv1.ClusterRoleBinding) (*rbacv1.ClusterRoleBinding, error) {
			crb.RoleRef = rbacv1.RoleRef{
				Name:     "system:auth-delegator",
				Kind:     "ClusterRole",
				APIGroup: rbacv1.GroupName,
			}
			crb.Subjects = []rbacv1.Subject{
				{
					Kind:      rbacv1.ServiceAccountKind,
					Name:      resources.MetricsServerServiceAccountName,
					Namespace: metav1.NamespaceSystem,
				},
			}
			return crb, nil
		}
	}
}
