[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

LM Operator is a [K8s Operator](https://coreos.com/operators/) that provides a K8s API to manage the installation of [LM](http://servicelifecyclemanager.com/2.1.0/) deployments.

# Limitations

* assumes that dependant services Cassandra, Elasticsearch, Kafka and Vault have already been installed, for example by using [Helm Foundation](http://servicelifecyclemanager.com/2.1.0/installation/lm/production/install-lm/).

# Developing

- [Developing LM Operator](./docs/developing.md) - docs for developers to build and install the driver
- [Installing LM Operator](./docs/installation.md) - installing LM Operator
- [Using LM Operator](./docs/using.md) - how to use the LM Operator
