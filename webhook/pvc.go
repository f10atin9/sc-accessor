package webhook

import (
	"context"
	"fmt"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"storageclass-accessor/client/apis/accessor/v1alpha1"
)

var reviewResponse = &admissionv1.AdmissionResponse{
	Allowed: true,
	Result:  &metav1.Status{},
}

func admitPVC(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	klog.V(2).Info("admitting pvc")

	if !(ar.Request.Operation == admissionv1.Update || ar.Request.Operation == admissionv1.Create) {
		return reviewResponse
	}

	raw := ar.Request.Object.Raw
	//oldRaw := ar.Request.OldObject.Raw

	deserializer := codecs.UniversalDeserializer()
	pvc := &corev1.PersistentVolumeClaim{}
	if _, _, err := deserializer.Decode(raw, nil, pvc); err != nil {
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}
	return decidePVCV1(pvc)

}

func decidePVCV1(pvc *corev1.PersistentVolumeClaim) *admissionv1.AdmissionResponse {
	if err := ValidateV1PVC(pvc); err != nil {
		reviewResponse.Allowed = false
		reviewResponse.Result.Message = err.Error()
	}
	return reviewResponse
}

func ValidateV1PVC(pvc *corev1.PersistentVolumeClaim) error {
	storageClassName := *pvc.Spec.StorageClassName
	namespace := pvc.Namespace
	// TODO: GET CR by accessor Name and validate
	cfg, err := config.GetConfig()
	if err != nil {
		return err
	}

	var cli client.Client
	opts := client.Options{}
	scheme := runtime.NewScheme()
	_ = v1alpha1.AddToScheme(scheme)
	opts.Scheme = scheme
	cli, err = client.New(cfg, opts)
	if err != nil {
		return err
	}

	accessor := &v1alpha1.Accessor{}
	err = cli.Get(context.Background(), types.NamespacedName{Namespace: "", Name: storageClassName + "-accessor"}, accessor)
	if err != nil {
		return err
	}
	for _, allowNS := range accessor.Spec.AllowedNamespace {
		if allowNS == namespace {
			return nil
		}
	}
	return fmt.Errorf("The storageClass: %s not allowed create pvc in the namespace: %s ", storageClassName, namespace)
}
