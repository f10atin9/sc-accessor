# Validating Webhook

The snapshot validating webhook is an HTTP callback which responds to [admission requests](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/). It is part of a larger [plan](https://github.com/kubernetes/enhancements/blob/master/keps/sig-storage/177-volume-snapshot/tighten-validation-webhook-crd.md) to tighten validation for volume snapshot objects. This webhook introduces the [ratcheting validation](https://github.com/kubernetes/enhancements/blob/master/keps/sig-storage/177-volume-snapshot/tighten-validation-webhook-crd.md#backwards-compatibility) mechanism targeting the tighter validation. The cluster admin or Kubernetes distribution admin should install the webhook alongside the snapshot controllers and CRDs.

> :warning: **WARNING**: Cluster admins choosing not to install the webhook server and participate in the phased release process can cause future problems when upgrading from `v1beta1` to `v1` volumesnapshot API, if there are currently persisted objects which fail the new stricter validation. Potential impacts include being unable to delete invalid snapshot objects.

## Prerequisites

The following are prerequisites to use this validating webhook:

- K8s version 1.17+ (v1.9+ to use `admissionregistration.k8s.io/v1beta1`, v1.16+ to use `admissionregistration.k8s.io/v1`, v1.17+ to use  `snapshot.storage.k8s.io/v1beta1`)
- ValidatingAdmissionWebhook is [enabled](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/#prerequisites). (in v1.18+ it will be enabled by default)
- API `admissionregistration.k8s.io/v1beta1` or `admissionregistration.k8s.io/v1` is enabled.

## How to build the webhook

Build the binary

```bash
make 
```

Build the docker image

```bash
docker build --network host -t f10atin9/webhook-accessor:v1.0 .
```


#### Method

These commands should be run from the top level directory.

1. Run the `create-cert.sh` script. Note using the default namespace will allow anyone with access to that namespace to read your secret. It is recommended to change the namespace in all the files and the commands given below.


    ```bash
    # This script will create a TLS certificate signed by the [cluster](https://kubernetes.io/docs/tasks/tls/managing-tls-in-a-cluster/). It will place the public and private key into a secret on the cluster.
    ./deploy/create-cert.sh --service pvc-accessor-service --secret pvc-validation-secret --namespace default # Make sure to use a different namespace
    ```

2. Patch the `ValidatingWebhookConfiguration` file from the template, filling in the CA bundle field.

    ```bash
    cat ./deploy/admission-configuration-template | ./deploy/patch-ca-bundle.sh > ./deploy/admission-configuration.yaml
    ```

3. Change the namespace in the generated `admission-configuration.yaml` file. Change the namespace in the service and deployment in the `webhook.yaml` file.

4. Create the deployment, service and admission configuration objects on the cluster.

    ```bash
    kubectl apply -f ./deploy
    ```
