# Using LM Operator

## Install LM

Create an [LM](http://servicelifecyclemanager.com/2.1.0/) deployment by running:

```
kubectl apply -f deploy/crds/alm.yaml
```

where alm.yaml looks like this:

```
apiVersion: com.accantosystems.stratoss/v1alpha1
kind: ALM
metadata:
  name: awesome
spec:
  # install LM in secure mode or not?
  secure: true
  # additional Spring profiles to configure for each LM service
  springProfilesActive: ""
  # deploymentType determines the number of LM service replicas to deploy. It also sets the CPU/memory
  # resource request/limits and the Java heap size.
  # Allowed values: tiny|middle|ha
  deploymentType: Middle
  # Docker registry from which to fetch LM service images
  dockerRepo: 10.220.217.248:32736
  # URL of LM release descriptor, detailing which versions of LM services to install
  release: http://10.220.217.248:8086/accanto/lm-operator-releases/raw/master/2.1.0-SNAPSHOT.yaml
  configurator:
    # run LM configurator or not
    Run: true
  conductor:
    JVMOptions: -Xmx256m
  brent:
    JVMOptions: -Xmx256m
  apollo:
    JVMOptions: -Xmx256m
  galileo:
    JVMOptions: -Xmx1024m
  talledega:
    JVMOptions: -Xmx1024m
  daytona:
    JVMOptions: -Xmx1024m
  nimrod:
    JVMOptions: -Xmx256m
  ishtar:
    JVMOptions: -Xmx256m
  relay:
    JVMOptions: -Xmx256m
  watchtower:
    JVMOptions: -Xmx1024m
  doki:
    JVMOptions: -Xmx512m
```

### LM Release Descriptor

An LM release descriptor defines which version of each LM microservice to install. For example:

```
configurator:
  version: 2.1.0-alpha-233
conductor:
  version: 2.1.0-alpha-233
brent:
  version: 2.1.0-alpha-233
apollo:
  version: 2.1.0-alpha-233
galileo:
  version: 2.1.0-alpha-233
talledega:
  version: 2.1.0-alpha-233
daytona:
  version: 2.1.0-alpha-233
nimrod:
  version: 2.1.0-alpha-233
ishtar:
  version: 2.1.0-alpha-233
relay:
  version: 2.1.0-alpha-233
watchtower:
  version: 2.1.0-alpha-233
doki:
  version: 2.1.0-alpha-233
```

## Uninstall LM

Uninstall LM by deleting the ALM instance:

```
kubectl delete ALM middle-alm
```

This will remove all LM services. It may leave some K8s artifacts (secrets) lying around, which you will have to clean up manually:

```

kubectl delete secret lm-certs lm-client-credentials lm-keystore  nimrod-tls  brent-tls  ishtar-tls  
```
