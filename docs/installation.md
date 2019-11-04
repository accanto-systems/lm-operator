# LM Operator Installation

Run the following commands to install the LM Operator in to an existing Kubernetes cluster:

```
kubectl apply -f deploy/service_account.yaml
kubectl apply -f deploy/role.yaml
kubectl apply -f deploy/role_binding.yaml
kubectl apply -f deploy/crds/com_v1alpha1_alm_crd.yaml
kubectl apply -f deploy/operator.yaml
```