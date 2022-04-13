# Prometheus examples

### Setup the basic Prometheus example

These are the steps which will be executed by the script:

* Connect to the IBM Cloud cluster
* Install OLM and Prometheus operator
* Install the example application
* Setup for of Prometheus for the example Application
* Verify the example application

```sh
sh setup-basic-prometheus-operator.sh
```

_Note:_ The documentation related to that example you find [here](https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/user-guides/getting-started.md).

### Setup FrontendOperator example

These are the steps which will be executed by the script:

* Connect to the IBM Cloud cluster
* Setup the FrontendOperator 
* Reconfigure the Prometheus instance

```sh
setup-frontendoperator-example-monitoring.sh
```

_Note:_ The documentation related to that example you find [here](https://wp.me/paelj4-1iv).

