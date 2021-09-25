package webhook

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"storageclass-accessor/client/apis/accessor/v1alpha1"
)

func validateNameSpace(resource, reqName, reqNameSpace string, accessor *v1alpha1.Accessor) error {
	klog.Info("start validate namespace")
	for _, allowedNS := range accessor.Spec.AllowedNamespace {
		if allowedNS == reqNameSpace {
			return nil
		}
	}
	klog.Error(fmt.Sprintf("%s %s don't allowed create in the namespace: %s", resource, reqName, reqNameSpace))
	return fmt.Errorf("The storageClass: %s not allowed create %s in the namespace: %s ", accessor.Spec.StorageClass, resource, reqNameSpace)
}

func validateWorkSpace(resource, reqName, reqNameSpace string, accessor *v1alpha1.Accessor) error {
	klog.Info("start validate workspace")
	cli, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		klog.Error("init a client failed, err:", err)
		return err
	}
	ns := &corev1.Namespace{}
	err = cli.Get(context.Background(), types.NamespacedName{Namespace: "", Name: reqNameSpace}, ns)
	if err != nil {
		klog.Error("client get namespace failed, err:", err)
		return err
	}
	var reqWorkSpace string
	var exist bool
	if reqWorkSpace, exist = ns.Labels["kubesphere.io/workspace"]; !exist {
		klog.Error("Can't get the workspace from the namespace " + ns.Name)
		return err
	}
	for _, allowWorkSpace := range accessor.Spec.AllowedWorkspace {
		if reqWorkSpace == allowWorkSpace {
			return nil
		}
	}
	klog.Error(fmt.Sprintf("%s %s don't allowed create in the workspace: %s", resource, reqName, reqWorkSpace))
	return fmt.Errorf("The storageClass: %s not allowed create %s in the workspace: %s ", accessor.Spec.StorageClass, resource, reqWorkSpace)
}
