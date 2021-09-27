package webhook

import (
	"context"
	admissionv1 "k8s.io/api/admission/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"

	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v4/apis/volumesnapshot/v1"
)

func admitSnapshot(ar admissionv1.AdmissionReview) *admissionv1.AdmissionResponse {
	klog.Info("admitting snapshot")

	if !(ar.Request.Operation == admissionv1.Delete || ar.Request.Operation == admissionv1.Create) {
		return reviewResponse
	}
	raw := ar.Request.Object.Raw

	var newSnapshot *snapshotv1.VolumeSnapshot

	switch ar.Request.Operation {
	case admissionv1.Create:
		deserializer := codecs.UniversalDeserializer()
		snapshot := &snapshotv1.VolumeSnapshot{}
		obj, _, err := deserializer.Decode(raw, nil, snapshot)
		if err != nil {
			klog.Error(err)
			return toV1AdmissionResponse(err)
		}
		var ok bool
		newSnapshot, ok = obj.(*snapshotv1.VolumeSnapshot)
		if !ok {
			klog.Error("obj can't exchange to snapshot object")
			return toV1AdmissionResponse(err)
		}
	case admissionv1.Delete:
		snapshotInfo := types.NamespacedName{
			Namespace: ar.Request.Namespace,
			Name:      ar.Request.Name,
		}

		cfg, err := config.GetConfig()
		if err != nil {
			return toV1AdmissionResponse(err)
		}
		var cli client.Client
		opts := client.Options{}
		scheme := runtime.NewScheme()
		_ = snapshotv1.AddToScheme(scheme)
		opts.Scheme = scheme
		cli, err = client.New(cfg, opts)
		if err != nil {
			return toV1AdmissionResponse(err)
		}
		snapshot := &snapshotv1.VolumeSnapshot{}
		err = cli.Get(context.Background(), snapshotInfo, snapshot)
		if err != nil {
			klog.Error("get target Delete Snapshot from client failed, err:", err)
			return toV1AdmissionResponse(err)
		}
		newSnapshot = snapshot
	}

	reqSnapshot := reqInfo{
		resource:         "volumeSnapshot",
		name:             newSnapshot.Name,
		namespace:        newSnapshot.Namespace,
		operator:         string(ar.Request.Operation),
		storageClassName: *newSnapshot.Spec.VolumeSnapshotClassName,
	}
	return decideSnapshot(reqSnapshot)
}

func decideSnapshot(snapshot reqInfo) *admissionv1.AdmissionResponse {
	accessor, err := getAccessor(snapshot.storageClassName)
	if err != nil {
		//TODO If not found , pass or not?
		return toV1AdmissionResponse(err)
	} else if accessor == nil {
		return reviewResponse
	}


	if err = validateNameSpace(snapshot, accessor); err != nil {
		return toV1AdmissionResponse(err)
	}

	if err = validateWorkSpace(snapshot, accessor); err != nil {
		return toV1AdmissionResponse(err)
	}
	return reviewResponse
}


//TODO getStorageClassByVolumeSnapshotClass

//func getStorageClassByVolumeSnapshotClass(snapshotClassName string) (storageclassName string, err error) {
//	var cli client.Client
//	opts := client.Options{}
//	scheme := runtime.NewScheme()
//	_ = snapshotv1.AddToScheme(scheme)
//	opts.Scheme = scheme
//	cli, err = client.New(config.GetConfigOrDie(), opts)
//	if err != nil {
//		return "", err
//	}
//	snapshotClass := &snapshotv1.VolumeSnapshotClass{}
//	err = cli.Get(context.Background(), types.NamespacedName{Namespace: "", Name: snapshotClassName}, snapshotClass)
//	if err != nil {
//		return "", err
//	}
//	return snapshotClass.
//}
