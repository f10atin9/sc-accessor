# storageclass-accessor
Help storageClass to manage cluster resource.

## install
### 1.install CRD and CR
```shell
kubectl apply -f  client/config/crds
kubectl apply -f example
```

### 2.create cert and secret
```bash
# This script will create a TLS certificate signed by the [cluster]It will place the public and private key into a secret on the cluster.
./deploy/create-cert.sh --service storageclass-accessor-service --secret accessor-validation-secret --namespace default # Make sure to use a different namespace
```
Move cert.pem and key.pem to the path "/etc/storageclass-accessor-webhook/certs"


### 3.Patch the `ValidatingWebhookConfiguration` file from the template, filling in the CA bundle field.
```shell
cat ./deploy/pvc-accessor-configuration-template | ./deploy/patch-ca-bundle.sh > ./deploy/pvc-accessor-configuration.yaml

```

### 4.build docker images
```shell
docker build --network host -t f10atin9/storageclass-accessor:v1.0 .
```

### 5.deploy 
```shell
kubectl apply -f deploy
```