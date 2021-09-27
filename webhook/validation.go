package webhook

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"storageclass-accessor/client/apis/accessor/v1alpha1"
)

func validateNameSpace(reqResource reqInfo, accessor *v1alpha1.Accessor) error {
	klog.Info("start validate namespace")

	//If not set, all namespace are allowed by default
	if len(accessor.Spec.AllowedNamespace) == 0 {
		return nil
	}

	for _, allowedNS := range accessor.Spec.AllowedNamespace {
		if allowedNS == reqResource.namespace {
			return nil
		}
	}
	klog.Error(fmt.Sprintf("%s %s does not allowed %s in the namespace: %s", reqResource.resource, reqResource.name, reqResource.operator, reqResource.namespace))
	return fmt.Errorf("The storageClass: %s does not allowed %s %s %s in the namespace: %s ", reqResource.storageClassName, reqResource.operator, reqResource.resource, reqResource.name, reqResource.namespace)
}

func validateWorkSpace(reqResource reqInfo, accessor *v1alpha1.Accessor) error {
	klog.Info("start validate workspace")

	//If not set, all workspace are allowed by default
	if len(accessor.Spec.AllowedWorkspace) == 0 {
		return nil
	}
	cli, err := client.New(config.GetConfigOrDie(), client.Options{})
	if err != nil {
		klog.Error("init a client failed, err:", err)
		return err
	}
	ns := &corev1.Namespace{}
	err = cli.Get(context.Background(), types.NamespacedName{Namespace: "", Name: reqResource.namespace}, ns)
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
	klog.Error(fmt.Sprintf("%s %s does not allowed %s in the workspace: %s", reqResource.resource, reqResource.name, reqResource.operator, reqWorkSpace))
	return fmt.Errorf("The storageClass: %s does not allowed %s %s %s in the workspace: %s ", reqResource.storageClassName, reqResource.operator, reqResource.resource, reqResource.name, reqWorkSpace)
}

func getAccessor(storageClassName string) (*v1alpha1.Accessor, error) {
	// get config
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	var cli client.Client
	opts := client.Options{}
	scheme := runtime.NewScheme()
	_ = v1alpha1.AddToScheme(scheme)
	opts.Scheme = scheme
	cli, err = client.New(cfg, opts)
	if err != nil {
		return nil, err
	}
	accessor := &v1alpha1.Accessor{}

	err = cli.Get(context.Background(), types.NamespacedName{Namespace: "", Name: storageClassName + "-accessor"}, accessor)
	if err != nil {
		//TODO If not found , pass or not?
		return nil, err
	}
	return accessor, nil
}
