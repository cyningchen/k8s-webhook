package main

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog"
)

func admissionRequired(ignoredList []string, metadata *metav1.ObjectMeta) bool {
	// skip special kubernetes system namespaces
	for _, namespace := range ignoredList {
		if metadata.Namespace == namespace {
			klog.V(2).Info("Skip validation for %v for it's in special namespace:%v", metadata.Name, metadata.Namespace)
			return false
		}
	}
	return true
}

func validationRequired(ignoredList []string, metadata *metav1.ObjectMeta) bool  {
	required := admissionRequired(ignoredList, metadata)
	klog.V(2).Info("Validation policy for %v/%v: required:%v", metadata.Namespace, metadata.Name, required)
	return required
}
