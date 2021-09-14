package webhook

import (
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/klog/v2"
)

var (
	PVCV1GVR = metav1.GroupVersionResource{Group: "storage.k8s.io", Version: "v1", Resource: "persistentvolumeclaims"}
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
	isUpdate := ar.Request.Operation == admissionv1.Update

	raw := ar.Request.Object.Raw
	oldRaw := ar.Request.OldObject.Raw

	deserializer := codecs.UniversalDeserializer()
	//switch ar.Request.Resource {
	//case PVCV1GVR:
	pvc := &corev1.PersistentVolumeClaim{}
	if _, _, err := deserializer.Decode(raw, nil, pvc); err != nil {
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}
	oldpvc := &corev1.PersistentVolumeClaim{}
	if _, _, err := deserializer.Decode(oldRaw, nil, pvc); err != nil {
		klog.Error(err)
		return toV1AdmissionResponse(err)
	}
	return decidePVCV1(pvc, oldpvc, isUpdate)
	//storageClass := pvc.Spec.StorageClassName
	//namespace := pvc.Namespace
	//if err := ValidateV1PVC(*storageClass, namespace); err != nil {
	//	toV1AdmissionResponse(err)
	//}
	//}
}

func decidePVCV1(pvc, oldpvc *corev1.PersistentVolumeClaim, isUpdate bool) *admissionv1.AdmissionResponse {
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
