package main

import (
	"context"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"testing"
)

func Test_main(t *testing.T) {
	cfg, err := config.GetConfig()
	if err != nil {
		fmt.Println("getConfig err:", err)
	}

	var cli client.Client
	cli, err = client.New(cfg, client.Options{})
	if err != nil {
		fmt.Println("init client failed:", err)
	}
	pvc := &corev1.PersistentVolumeClaim{}
	err = cli.Get(context.Background(), types.NamespacedName{
		Namespace: "default",
		Name:      "tt1",
	}, pvc)
	if err != nil {
		fmt.Println("get pvc tt1 failed:", err)
	} else {
		fmt.Println("get pvc tt1 successed")
	}

}
