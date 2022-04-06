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

function setupFrontendOperator () {
   # deploy the custom operator and an example application deployed by the custom operator
   export REGISTRY='quay.io'
   export ORG='tsuedbroecker'
   export CONTROLLER_IMAGE='frontendcontroller-monitoring:v0.0.1'
   cd ..  
   make deploy IMG="$REGISTRY/$ORG/$CONTROLLER_IMAGE"
   kubectl apply -f config/samples/multitenancy_v1alpha1_tenancyfrontend.yaml -n default
   cd ./prometheus-example
   kubectl get servicemonitor -n frontendoperator-system
}

function reconfigurePrometheusInstance () {
    # we reconfigure our Prometheus instance with 'wildcards' for search after servicemonitors
    echo "9-prom-instance-servicemonitor-selector-reconfigure.yaml"
    kubectl apply -f 9-prom-instance-servicemonitor-selector-reconfigure.yaml -n default
    kubectl get prometheus -n default
    kubectl get statefulset -n default

    # inspect the prometheus pod 
    # ===============
    # kubectl exec -it prometheus-prometheus-0 -n default -- /bin/sh 
    # --------------------------------------------------------------
    # inside the pod: 
    # ps
    # cd ..
    # cd /etc/prometheus/config_out
    # ls 
    # cat prometheus.env.yaml 
    # cat prometheus.env.yaml | grep -i -A 10 "job_name:"
    # cat /etc/prometheus/config_out/prometheus.env.yaml | grep -i -A 10 "job_name: default/prometheus-prometheus/0"
}

# **********************************************************************************
# Execution
# **********************************************************************************

echo "************************************"
echo " connectToIBMCluster"
echo "************************************"

connectToIBMCluster

echo "************************************"
echo " setupFrontendOperator "
echo "************************************"

setupFrontendOperator 

echo "************************************"
echo "reconfigurePrometheusInstance "
echo "************************************"

reconfigurePrometheusInstance