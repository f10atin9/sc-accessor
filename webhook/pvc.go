package webhook

import (
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

func admitPVC(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	klog.V(2).Info("admitting pvc")

	reviewResponse := &admissionv1.AdmissionResponse{
		Allowed: true,
		Result:  &metav1.Status{},
	}
	if !(ar.Request.Operation == admissionv1.Update || ar.Request.Operation == admissionv1.Create) {
		return reviewResponse
	}
	//isUpdate := ar.Request.Operation == admissionv1.Update

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
	reviewResponse := &admissionv1.AdmissionResponse{
		Allowed: true,
		Result:  &metav1.Status{},
	}

	if err := ValidateV1PVC(pvc); err != nil {
		reviewResponse.Allowed = false
		reviewResponse.Result.Message = err.Error()
	}
	return reviewResponse
}
