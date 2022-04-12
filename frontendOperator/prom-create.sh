#!/bin/bash

# OLM
operator-sdk olm install latest
echo "Press any key to move on"
read input

# Prometheus Operator
kubectl create -f https://operatorhub.io/install/prometheus.yaml
echo "Press any key to move on"
read input

# Prometheus Configuration
kubectl create namespace monitoring
kubectl create -f prom-serviceaccount.yaml -n monitoring
kubectl create -f prom-clusterrole.yaml -n monitoring
kubectl create -f prom-clusterrolebinding.yaml -n monitoring
kubectl create -f prom-instance.yaml -n monitoring
kubectl create -f prom-loadbalancer.yaml -n monitoring
echo "Press any key to move on"
read input

make generate
make manifests


# Example Operator



EXTERNAL_IP=$(kubectl get service prometheus -n monitoring | grep prometheus |  awk '{print $4;}')
PORT=$(kubectl get service prometheus -n monitoring | grep prometheus |  awk '{print $5;}'| sed 's/\(.*\):.*/\1/g')
echo "http://$EXTERNAL_IP:$PORT"

# Build
make generate
make manifests

# Build container
make docker-build IMG="$REGISTRY/$ORG/$CONTROLLER_IMAGE"
# Push container
podman login quay.io
podman push "$REGISTRY/$ORG/$CONTROLLER_IMAGE"

# Deploy
export REGISTRY='quay.io'
export ORG='tsuedbroecker'
export CONTROLLER_IMAGE='frontendcontroller-monitoring:v0.0.1'
make deploy IMG="$REGISTRY/$ORG/$CONTROLLER_IMAGE"
kubectl apply -f config/samples/multitenancy_v1alpha1_tenancyfrontend.yaml -n default
