package main

import (
	"context"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"webhook/api/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/persistentvolumeclaims,mutating=false,failurePolicy=fail,groups="",resources=persistentvolumeclaims,verbs=create;update,versions=v1,name=pvc.ks.io
type pvcValidator struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (v *pvcValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	webhookLog.Info("starting webHook handle")
	pvc := &corev1.PersistentVolumeClaim{}

	err := v.decoder.Decode(req, pvc)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	accessor := &v1alpha1.StorageClassAccess{}
	accessorName := *pvc.Spec.StorageClassName + "-accessor"
	key := types.NamespacedName{
		Namespace: "",
		Name:      accessorName,
	}
	if err := v.Client.Get(ctx, key, accessor); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}
	for _, allowNS := range accessor.Spec.AllowedNamespace {
		if allowNS == pvc.Namespace {
			return admission.Allowed("")
		}
	}
	return admission.Denied("this storageClass is not allowed creating pvc in" + pvc.Namespace)
}
