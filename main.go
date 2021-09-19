package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/runtime"
	"log"
	"net/http"
	"storageclass-accessor/accessor/v1alpha1"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

func main() {
	http.HandleFunc("/pvc", pvcHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func pvcHandler(w http.ResponseWriter, r *http.Request) {
	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		log.Println("read request body failed")
		return
	}

	pvc := &corev1.PersistentVolumeClaim{}
	err = json.Unmarshal(reqBody, pvc)
	if err != nil {
		log.Println("unmarshal json failed")
		return
	}

	cfg, err := config.GetConfig()
	if err != nil {
		log.Println("getConfig err:", err)
	}

	var cli1 client.Client
	opts := client.Options{}
	scheme := runtime.NewScheme()
	_ = v1alpha1.AddToScheme(scheme)
	opts.Scheme = scheme
	cli1, err = client.New(cfg, opts)

	storageClass := *pvc.Spec.StorageClassName
	namespace := pvc.Namespace

	accessor := &v1alpha1.Accessor{}

	err = cli1.Get(
		context.Background(),
		types.NamespacedName{Namespace: "", Name: storageClass + "-accessor"},
		accessor,
	)
	if err != nil {
		log.Println("get accessor failed:", err)
		return
	}
	fmt.Println("hello")
	for _, ele := range accessor.Spec.AllowedNamespace {
		fmt.Println("ele is :", ele)
		if ele == namespace {
			var pvcClient client.Client
			pvcClient, _ = client.New(config.GetConfigOrDie(), client.Options{})
			err := pvcClient.Create(context.Background(), pvc)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)

				log.Fatal("post req to sc create pvc failed:", err)

				return
			}
			w.WriteHeader(http.StatusOK)
			return
		}
	}
	w.WriteHeader(http.StatusBadRequest)
	return
}
