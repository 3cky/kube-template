// Copyright © 2015-2018 Victor Antonovich <victor@antonovich.me>
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

// Code generated by "go:generate go run cmd/gentemplate/main.go". DO NOT EDIT.

package main

import (
	corev1 "k8s.io/api/core/v1"
)

func kubeObjectsFuncMap(dm *DependencyManager) map[string]interface{} {
	return map[string]interface{}{
		"pods":                   pods(dm),
		"services":               services(dm),
		"replicationcontrollers": replicationcontrollers(dm),
		"events":                 events(dm),
		"endpoints":              endpoints(dm),
		"nodes":                  nodes(dm),
		"namespaces":             namespaces(dm),
		"componentstatuses":      componentstatuses(dm),
		"configmaps":             configmaps(dm),
		"limitranges":            limitranges(dm),
		"persistentvolumes":      persistentvolumes(dm),
		"persistentvolumeclaims": persistentvolumeclaims(dm),
		"podtemplates":           podtemplates(dm),
		"resourcequotas":         resourcequotas(dm),
		"secrets":                secrets(dm),
		"serviceaccounts":        serviceaccounts(dm),
	}
}

// {{pods "selector" "namespace"}}
func pods(dm *DependencyManager) func(...string) ([]corev1.Pod, error) {
	return func(s ...string) ([]corev1.Pod, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.Pods(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{services "selector" "namespace"}}
func services(dm *DependencyManager) func(...string) ([]corev1.Service, error) {
	return func(s ...string) ([]corev1.Service, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.Services(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{replicationcontrollers "selector" "namespace"}}
func replicationcontrollers(dm *DependencyManager) func(...string) ([]corev1.ReplicationController, error) {
	return func(s ...string) ([]corev1.ReplicationController, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.ReplicationControllers(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{events "selector" "namespace"}}
func events(dm *DependencyManager) func(...string) ([]corev1.Event, error) {
	return func(s ...string) ([]corev1.Event, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.Events(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{endpoints "selector" "namespace"}}
func endpoints(dm *DependencyManager) func(...string) ([]corev1.Endpoints, error) {
	return func(s ...string) ([]corev1.Endpoints, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.Endpoints(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{nodes "selector"}}
func nodes(dm *DependencyManager) func(...string) ([]corev1.Node, error) {
	return func(s ...string) ([]corev1.Node, error) {
		if selector, err := parseSelector(s...); err == nil {
			return dm.Nodes(selector)
		} else {
			return nil, err
		}
	}
}

// {{namespaces "selector"}}
func namespaces(dm *DependencyManager) func(...string) ([]corev1.Namespace, error) {
	return func(s ...string) ([]corev1.Namespace, error) {
		if selector, err := parseSelector(s...); err == nil {
			return dm.Namespaces(selector)
		} else {
			return nil, err
		}
	}
}

// {{componentstatuses "selector"}}
func componentstatuses(dm *DependencyManager) func(...string) ([]corev1.ComponentStatus, error) {
	return func(s ...string) ([]corev1.ComponentStatus, error) {
		if selector, err := parseSelector(s...); err == nil {
			return dm.ComponentStatuses(selector)
		} else {
			return nil, err
		}
	}
}

// {{configmaps "selector" "namespace"}}
func configmaps(dm *DependencyManager) func(...string) ([]corev1.ConfigMap, error) {
	return func(s ...string) ([]corev1.ConfigMap, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.ConfigMaps(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{limitranges "selector" "namespace"}}
func limitranges(dm *DependencyManager) func(...string) ([]corev1.LimitRange, error) {
	return func(s ...string) ([]corev1.LimitRange, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.LimitRanges(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{persistentvolumes "selector"}}
func persistentvolumes(dm *DependencyManager) func(...string) ([]corev1.PersistentVolume, error) {
	return func(s ...string) ([]corev1.PersistentVolume, error) {
		if selector, err := parseSelector(s...); err == nil {
			return dm.PersistentVolumes(selector)
		} else {
			return nil, err
		}
	}
}

// {{persistentvolumeclaims "selector" "namespace"}}
func persistentvolumeclaims(dm *DependencyManager) func(...string) ([]corev1.PersistentVolumeClaim, error) {
	return func(s ...string) ([]corev1.PersistentVolumeClaim, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.PersistentVolumeClaims(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{podtemplates "selector" "namespace"}}
func podtemplates(dm *DependencyManager) func(...string) ([]corev1.PodTemplate, error) {
	return func(s ...string) ([]corev1.PodTemplate, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.PodTemplates(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{resourcequotas "selector" "namespace"}}
func resourcequotas(dm *DependencyManager) func(...string) ([]corev1.ResourceQuota, error) {
	return func(s ...string) ([]corev1.ResourceQuota, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.ResourceQuotas(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{secrets "selector" "namespace"}}
func secrets(dm *DependencyManager) func(...string) ([]corev1.Secret, error) {
	return func(s ...string) ([]corev1.Secret, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.Secrets(namespace, selector)
		} else {
			return nil, err
		}
	}
}

// {{serviceaccounts "selector" "namespace"}}
func serviceaccounts(dm *DependencyManager) func(...string) ([]corev1.ServiceAccount, error) {
	return func(s ...string) ([]corev1.ServiceAccount, error) {
		if namespace, selector, err := parseNamespaceSelector(s...); err == nil {
			return dm.ServiceAccounts(namespace, selector)
		} else {
			return nil, err
		}
	}
}
