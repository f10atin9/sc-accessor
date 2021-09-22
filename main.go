package main

import (
	"context"
	"encoding/json"
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

	// Listening create pvc request
	log.Println("start listening")

	http.HandleFunc("*/pvc", pvcCreateHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func pvcCreateHandler(w http.ResponseWriter, r *http.Request) {
	// read request
	reqBody, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	var msg string
	if err != nil {
		msg = "read request body failed"
		log.Println(msg)
		w.Write([]byte(msg))
		return
	}

	//get pvc object
	pvc := &corev1.PersistentVolumeClaim{}
	err = json.Unmarshal(reqBody, pvc)
	if err != nil {
		msg = "unmarshal json failed"
		log.Println(msg)
		w.Write([]byte("unmarshal json failed"))
		return
	}
	storageClass := *pvc.Spec.StorageClassName
	namespace := pvc.Namespace

	//init accessor client
	cfg, err := config.GetConfig()
	if err != nil {
		msg = "getConfig err:" + err.Error()
		log.Println(msg)
		w.Write([]byte(msg))
	}
	var cli1 client.Client
	opts := client.Options{}
	scheme := runtime.NewScheme()
	_ = v1alpha1.AddToScheme(scheme)
	opts.Scheme = scheme
	cli1, err = client.New(cfg, opts)

	//get accessor by pvc storageClassName
	accessor := &v1alpha1.Accessor{}
	err = cli1.Get(
		context.Background(),
		types.NamespacedName{Namespace: "", Name: storageClass + "-accessor"},
		accessor,
	)
	if err != nil {
		msg = "get accessor failed:" + err.Error()
		log.Println(msg)
		w.Write([]byte(msg))
		return
	}

	for _, ns := range accessor.Spec.AllowedNamespace {
		// if in allowNS, create pvc
		if ns == namespace {
			var pvcClient client.Client
			pvcClient, _ = client.New(config.GetConfigOrDie(), client.Options{})
			err := pvcClient.Create(context.Background(), pvc)
			if err != nil {
				w.WriteHeader(http.StatusForbidden)
				msg = "post req to sc create pvc failed:" + err.Error()
				log.Fatal(msg)
				w.Write([]byte(msg))
				return
			}
			w.WriteHeader(http.StatusOK)
			msg = storageClass + " pvc \"" + pvc.Name + "\"created in the namespace " + namespace
			w.Write([]byte("pvc"))
			return
		}
	}

	// else reject request
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte("It is not allowed to create " + storageClass + " pvc in the namespace" + namespace))
	return
}
