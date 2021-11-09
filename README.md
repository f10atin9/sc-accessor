# storageclass-accessor
***
## 
The storageclass-accessor webhook is an HTTP callback which responds to admission requests.
When creating and deleting the PVC, it will take out the accessor related to this storage class, and the request will be allowed only when all accessors pass the verification.
Users can create accessors and set `namespaceSelector` to achieve **namespace-level** management on the storage class which provisions PVC.

## Quick Start
***
The guide describes how to deploy a storage class accessor webhook to a cluster and provides an example accessor based on csi-qingcloud.
### 1. Install CRD and CR
```shell
kubectl create -f  client/config/crds
```

### 2. Create certificate and secret
```bash
# This script will create a TLS certificate signed by the [cluster]It will place the public and private key into a secret on the cluster.
./deploy/create-cert.sh --service storageclass-accessor-service --secret accessor-validation-secret --namespace default # Make sure to use a different namespace
```
Move `cert.pem` and `key.pem` to the path `/etc/storageclass-accessor-webhook/certs`.


### 3. Patch the `ValidatingWebhookConfiguration` file from the template and fill in the CA bundle field
```shell
cat ./deploy/pvc-accessor-configuration-template | ./deploy/patch-ca-bundle.sh > ./deploy/pvc-accessor-configuration.yaml
```

### 4. Build Docker images
```shell
docker build --network host -t f10atin9/storageclass-accessor:v1.0 .
```

### 5. Deploy 
```shell
kubectl apply -f deploy
```

### 6. Apply example CR
```shell
kubectl apply -f example
```

## Next
***
If you need to customize the accessor rules, write the YAML file according to the example YAML and apply it to your cluster. The accessor rules will be implemented when a storage class is requested.

## Example
***
This example YAML shows how to set the accessor, and you can define your own `namespaceSelector` according to your needs.
```yaml
apiVersion: storage.kubesphere.io/v1alpha1
kind: Accessor
metadata:
  name: csi-qingcloud-accessor
spec:
  storageClassName: "csi-qingcloud"
  namespaceSelector:
    fieldSelector:
      - fieldExpressions:
          - field: "Name"
            operator: "In"
            values: ["default"]
    labelSelector:
      - matchExpressions:
          - key: "app"
            operator: "In"
            values: ["test-app"]
          - key: "role"
            operator: "In"
            values: ["admin", "user"]
      - matchExpressions:
          - key: "app"
            operator: "In"
            values: ["test-app2"]
```

- When there are multiple rules in `fieldExpressions` or `matchExpressions`, all the rules need to pass the verification before they can be implemented.
- If there are multiple `fieldExpressions`, only one of them needs to pass the verification. So does `matchExpressions`.
- When both the `fieldSelector` and `labelSelector` pass the verification, the `namespaceSelector` is judged to pass the verification as well.
- If a storage class is mentioned by multiple accessors, it needs to pass all the rules of these accessors.