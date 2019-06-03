package main

import (
	"k8s.io/api/admission/v1beta1"
	"k8s.io/klog"
	"k8s.io/api/apps/v1"
	v12 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

const (
	IgnorePrefix = "hadoop"
)

var (
	ignoredNamespaces = []string{
		metav1.NamespaceSystem,
		metav1.NamespacePublic,
	}
)


func validateLabels(vl v1beta1.AdmissionReview) *v1beta1.AdmissionResponse{
	var (
		availableLabels map[string]string
		objectMeta	*metav1.ObjectMeta
		resourceNamespace, resourceName string
	)
	klog.V(2).Info("validating labels")
	req := vl.Request
	switch req.Kind.Kind {
	case "Deployment":
		var deploy v1.Deployment
		deserializer := codecs.UniversalDeserializer()
		if _, _, err := deserializer.Decode(req.Object.Raw, nil, &deploy); err != nil {
			klog.Error(err)
			return toAdmissionResponse(err)
		}
		resourceName, resourceNamespace, objectMeta = deploy.Name, deploy.Namespace, &deploy.ObjectMeta
		availableLabels = deploy.Labels
	case "Service":
		var service v12.Service
		deserializer := codecs.UniversalDeserializer()
		if _, _, err := deserializer.Decode(req.Object.Raw, nil, &service); err != nil {
			klog.Error(err)
			return toAdmissionResponse(err)
		}
		resourceName, resourceNamespace, objectMeta = service.Name, service.Namespace, &service.ObjectMeta
		availableLabels = service.Labels
	case "Pod":
		var pod v12.Pod
		deserializer := codecs.UniversalDeserializer()
		if _, _, err := deserializer.Decode(req.Object.Raw, nil, &pod); err != nil {
			klog.Error(err)
			return toAdmissionResponse(err)
		}
		resourceName, resourceNamespace, objectMeta = pod.Name, pod.Namespace, &pod.ObjectMeta
		availableLabels = pod.Labels
	}
	if !validationRequired(ignoredNamespaces, objectMeta) {
		klog.V(2).Info("Skipping validation for %s/%s due to policy check", resourceNamespace, resourceName)
		return &v1beta1.AdmissionResponse{
			Allowed: true,
		}
	}
	allowed := false
	for k := range availableLabels{
		if strings.HasPrefix(k, IgnorePrefix){
			allowed = true
			break
		}
	}
	return  &v1beta1.AdmissionResponse{
		Allowed: allowed,
		Result: &metav1.Status{Message: strings.TrimSpace("required label are not set")},
	}
}