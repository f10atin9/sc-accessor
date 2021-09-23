package webhook

//import (
//	"context"
//	"fmt"
//	"storageclass-accessor/client/apis/accessor/v1alpha1"
//
//	corev1 "k8s.io/api/core/v1"
//	"sigs.k8s.io/controller-runtime/pkg/client"
//	"sigs.k8s.io/controller-runtime/pkg/client/config"
//)
//
//func ValidateV1PVC(pvc *corev1.PersistentVolumeClaim) error {
//	storageClassName := *pvc.Spec.StorageClassName
//	namespace := pvc.Namespace
//	// TODO: GET CR by accessor Name and validate
//	cfg, err := config.GetConfig()
//	if err != nil {
//		return err
//	}
//
//	var cli client.Client
//	cli, err = client.New(cfg, client.Options{})
//	if err != nil {
//		return err
//	}
//
//	accessList := &v1alpha1.AccessorList{}
//	opts := []client.ListOption{
//		client.MatchingFields{"allowedNamespace": namespace},
//	}
//	err = cli.List(context.Background(), accessList, opts...)
//	if err != nil {
//		return err
//	}
//	for _, ele := range accessList.Items {
//		if ele.Name == storageClassName {
//			return nil
//		}
//	}
//	return fmt.Errorf("The storageClass: %s not allowed create pvc in the namespace: %s ", storageClassName, namespace)
//}
