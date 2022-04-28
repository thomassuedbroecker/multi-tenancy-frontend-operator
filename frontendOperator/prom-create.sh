#!/bin/bash

# **************** Global variables

export REGISTRY='quay.io'
export ORG='tsuedbroecker'
export CONTROLLER_IMAGE='frontendcontroller-monitoring:v0.0.1'

# **********************************************************************************
# Functions
# **********************************************************************************

# OLM

function installOLM () {
    operator-sdk olm install latest
    echo "Press any key to move on"
    read input
}

# Prometheus Operator
function installPrometheusOperator () {
    kubectl create -f https://operatorhub.io/install/prometheus.yaml
    echo "Press any key to move on"
    read input
}

# Prometheus Configuration
function configurePrometheusOperator () {
    kubectl create namespace monitoring
    kubectl create -f prom-serviceaccount.yaml -n monitoring
    kubectl create -f prom-clusterrole.yaml -n monitoring
    kubectl create -f prom-clusterrolebinding.yaml -n monitoring
    kubectl create -f prom-instance.yaml -n monitoring
    kubectl create -f prom-loadbalancer.yaml -n monitoring

    kubectl get pods -n monitoring
    echo "Press any key to move on"
    read input
}

# Example Operator
function getPrometheusUIURL () {
    EXTERNAL_IP=$(kubectl get service prometheus -n monitoring | grep prometheus |  awk '{print $4;}')
    PORT=$(kubectl get service prometheus -n monitoring | grep prometheus |  awk '{print $5;}'| sed 's/\(.*\):.*/\1/g')
    echo "Access Prometheus UI: http://$EXTERNAL_IP:$PORT"
}

# Build
function buildAndUploadCustomOperator () {
  make generate
  make manifests
  # Build container
  make docker-build IMG="$REGISTRY/$ORG/$CONTROLLER_IMAGE"
  # Push container
  podman login quay.io
  podman push "$REGISTRY/$ORG/$CONTROLLER_IMAGE"
}

# Deploy
function deployCustomOperator () {
    make deploy IMG="$REGISTRY/$ORG/$CONTROLLER_IMAGE"
    kubectl apply -f config/samples/multitenancy_v1alpha1_tenancyfrontend.yaml -n default
}

# **********************************************************************************
# Execution
# **********************************************************************************

installOLM

installPrometheusOperator

configurePrometheusOperator

getPrometheusUIURL

# Uncomment the next line to use your own build
# ===============================
# buildAndUploadCustomOperator

deployCustomOperator

