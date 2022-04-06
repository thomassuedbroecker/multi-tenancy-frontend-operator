#!/bin/bash 

# **************** Global variables

 # ibm cloud specific
export RESOURCE_GROUP=default
export REGION=us-south
export CLUSTER_NAME=YOUR_CLUSTER


# **********************************************************************************
# Functions
# **********************************************************************************

function connectToIBMCluster () {

    ibmcloud target -g $RESOURCE_GROUP
    ibmcloud target -r $REGION

    ibmcloud ks cluster config --cluster $CLUSTER_NAME
    kubectl config current-context
}

function installOLMandPrometheus () {
    # install OLM
    operator-sdk olm install latest

    # install prometheus operator
    # Documention: https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/user-guides/getting-started.md
    # Installation of the operator: https://operatorhub.io/operator/prometheus
    # curl -sL https://github.com/operator-framework/operator-lifecycle-manager/releases/download/v0.20.0/install.sh | bash -s v0.20.0
    # Operator subscription need to have the operatorhubio-catalog entry in the olm namespace exist  
    kubectl get catalogsource -n olm | grep "operatorhubio-catalog"
    kubectl create -f https://operatorhub.io/install/prometheus.yaml  
}

function setupExampleApplication () {
    echo "1-app-example-deployment.yaml"
    kubectl apply -f 1-app-example-deployment.yaml -n default
    kubectl get deployments -n default

    echo "2-app-example-service.yaml"
    kubectl apply -f 2-app-service.yaml -n default
    kubectl get service -n default
}

function setupPrometheusForExampleApplication () {
    echo "3-prom-service-monitor.yaml"
    kubectl apply -f 3-prom-service-monitor.yaml -n default
    kubectl get servicemonitor -n default

    echo "4-prom-service-account.yaml"
    kubectl apply -f 4-prom-service-account.yaml -n default
    kubectl get serviceaccount -n default

    echo "5-prom-cluster-role.yaml"
    kubectl apply -f 5-prom-cluster-role.yaml -n default
    kubectl get clusterrole -n default

    echo "6-prom-cluster-role-binding.yaml"
    kubectl apply -f 6-prom-cluster-role-binding.yaml -n default

    echo "7-prom-instance-servicemonitor-selector-configured.yaml"
    kubectl apply -f 7-prom-instance-servicemonitor-selector-configured.yaml -n default

    echo "8-prom-expose-ui-loadbalancer.yaml"
    kubectl apply -f 8-prom-expose-ui-loadbalancer.yaml -n default
}

function verifyExampleApplication () {
    # verify the installation does the csv exist
    kubectl get csv -n default | grep "prometheusoperator"
    kubectl get pods -n default -o wide --show-labels
    kubectl get service -n default  -o wide --show-labels
    kubectl get servicemonitor -n default  -o wide --show-labels
    kubectl get prometheus -n default  -o wide --show-labels
    kubectl get prometheus prometheus -n default  -oyaml
    kubectl get clusterrole prometheus -n default  -o wide --show-labels
    kubectl get clusterrolebinding prometheus -n default  -o wide --show-labels
    kubectl get configmap -n default -o wide --show-labels
    kubectl get configmap prometheus-prometheus-rulefiles-0 -n default -oyaml
    kubectl get secret prometheus-prometheus -n default -o yaml
}

# **********************************************************************************
# Execution
# **********************************************************************************

connectToIBMCluster

setupExampleApplication

setupPrometheusForExampleApplication

verifyExampleApplication