# storageclass-accessor
***
## 
The storageclass-accessor webhook is an HTTP callback which responds to admission requests.
When creating and deleting the PVC, it will take out the accessor related to this storageclass, and the request will be allowed only when all accessors pass the verification.
Users can create accessor and set namespaceSelector to achieve **namespace-level** management on the StorageClass to create pvc

## Quick Start
***
The guide shows how to deploy StorageClass accessor webhook to the cluster. And provides an example accessor about csi-qingcloud.
### 1.install CRD and CR
```shell
kubectl create -f  client/config/crds
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

### 4.Deploy 
```shell
kubectl apply -f deploy
```

### 5.Write a CR


### 6.Apply CR
Use the `kubectl apply` command to make the accessor you created work.

### 7.Test
Now you can try to create a PVC. If it is created in a namespace that is not allowed, the following error will be output:

> Error from server: error when creating "PVC.yaml": admission webhook "pvc-accessor.storage.kubesphere.io" denied the request: The storageClass: **StorageClassName** does not allowed CREATE persistentVolumeClaim **PVC-NAME** in the namespace: **TARGET-NS**

## Accessor CR
***
A complete accessor should have the following fields:
 - spec.storageClassName

   The accessor knows the effective sc according to this field.
   - spec.namespaceSelector

     This field is used to fill in the limit of nameSpace, Including **labelSelector** and **fieldSelector**.
     - spec.namespaceSelector.labelSelector

       It's an **array of matchExpressions** . Manage whether nameSpace is available through the label of nameSpace.
          - spec.namespaceSelector.labelSelector.matchExpressions
     
          It's an **array of labelRule** . Every rule in the array needs to be verified.

          labelRule has the following fields:
     
     1.key
     
   
   - spec.namespaceSelector.fieldSelector

     Is an **array of fieldExpressions** .Manage whether nameSpace is available through the label of nameSpace.
   


## Example
***
The next few examples of yaml may be helpful for you to design Accessor:
### OnlyFieldSelector

#### Only one FieldExpression
```yaml
apiVersion: storage.kubesphere.io/v1alpha1
kind: Accessor
metadata:
  name: onlyFieldSelector-accessor
spec:
  storageClassName: "csi-qingcloud"
  namespaceSelector:
    fieldSelector:
      - fieldExpressions:
          - field: "Name"
            operator: "In"
            values: ["NS1"]
```
After applying this accessor, you can create the pvc of csi-qingcloud only in namespace.name which in this array :["NS1"]

***
More than one fieldExpressions are allowed in a fieldSelector.

And multiple rules are also allowed in fieldExpressions

#### Multiple FieldExpressions
```yaml
apiVersion: storage.kubesphere.io/v1alpha1
kind: Accessor
metadata:
  name: multipleFieldExpressions-accessor
spec:
  storageClassName: "csi-qingcloud"
  namespaceSelector:
    fieldSelector:
      - fieldExpressions:
          - field: "Name"
            operator: "In"
            values: ["NS1"]
      - fieldExpressions:
          - field: "Name"
            operator: "In"
            values: ["NS2", "NS3"]
```
You can create the pvc of csi-qingcloud in namespace which (nameSpace.Name in ["NS1"]) **or** (nameSpace.Name in ["NS2", "NS3"])

#### Multiple rule in one FieldExpressions
```yaml
apiVersion: storage.kubesphere.io/v1alpha1
kind: Accessor
metadata:
  name: multipleFieldExpressions-accessor
spec:
  storageClassName: "csi-qingcloud"
  namespaceSelector:
    fieldSelector:
      - fieldExpressions:
          - field: "Name"
            operator: "NotIn"
            values: ["NS1", "NS2"]
          - field: "Status"
            operator: "In"
            values: ["Active"]
```
You can create the pvc of csi-qingcloud only in namespace which (nameSpace.Name NotIn ["NS1", "NS2"]) **and** (nameSpace.Status in ["Active"])

It means that the rules in fieldExpressions must be followed at the same time.

### OnlyLabelSelector

####  Only one matchExpressions
```yaml
apiVersion: storage.kubesphere.io/v1alpha1
kind: Accessor
metadata:
  name: csi-qingcloud-accessor
spec:
  storageClassName: "csi-qingcloud"
  namespaceSelector:
    labelSelector:
      - matchExpressions:
          - key: "target-label"
            operator: "In"
            values: ["test-app"]
```
This requires nameSpace to have the key "app" tag and the value in this array: []

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

- 1. When there are multiple rules in a fieldExpressions or matchExpressions, all the rules need to pass the verification to pass.
- 2. If there are multiple fieldExpressions, only one of them needs to pass, and matchExpressions are the same.
- 3. When both the fieldSelector and labelSelector pass, the namespaceSelector is judged to pass.
- 4. If a StorageClass is mentioned by multiple accessors, it needs to pass all accessor rules.

## Notice
:warining: **Warning**:Too many accessors may cause unexpected errors in the webhook. It is recommended that one storageClass corresponds to one accessor.