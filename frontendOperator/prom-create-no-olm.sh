#!/bin/bash

# **************** Global variables

export REGISTRY='quay.io'
export ORG='tsuedbroecker'
export CONTROLLER_IMAGE='frontendcontroller-monitoring:v0.0.1'

# **********************************************************************************
# Functions
# **********************************************************************************

# Prometheus Operator
function installPrometheusOperator () {
    
    # kubectl create -f https://operatorhub.io/install/prometheus.yaml
    # curl https://raw.githubusercontent.com/prometheus-operator/prometheus-operator/main/bundle.yaml > prometheus-bundle.yaml
    kubectl create namespace monitoring
    kubectl create -f ./prom-bundle-monitoring.yaml
    kubectl get pods -n monitoring

    array=("prometheus-operator" )
  namespace=monitoring
  export STATUS_SUCCESS="Running"
  for i in "${array[@]}"
    do 
        echo ""
        echo "------------------------------------------------------------------------"
        echo "Check $i"
        while :
        do
            FIND=$i
            STATUS_CHECK=$(kubectl get pods -n $namespace | grep "$FIND" | awk '{print $3;}' | sed 's/"//g' | sed 's/,//g')
            echo "Status: $STATUS_CHECK"
            STATUS_VERIFICATION=$(echo "$STATUS_CHECK" | grep $STATUS_SUCCESS)
            if [ "$STATUS_VERIFICATION" = "$STATUS_SUCCESS" ]; then
                echo "$(date +'%F %H:%M:%S') Status: $FIND is Ready"
                echo "------------------------------------------------------------------------"
                break
            else
                echo "$(date +'%F %H:%M:%S') Status: $FIND($STATUS_CHECK)"
                echo "------------------------------------------------------------------------"
            fi
            sleep 3
        done
    done 
    
    echo "Press any key to move on"
    read input
}

# Prometheus Configuration
function configurePrometheusOperator () {

    kubectl create -f prom-instance.yaml -n monitoring
    kubectl create -f prom-service-nodeport.yaml -n monitoring
    kubectl create -f prom-serviceaccount.yaml -n monitoring
    kubectl create -f prom-clusterrole.yaml -n monitoring
    kubectl create -f prom-clusterrolebinding.yaml -n monitoring

    kubectl get pods -n monitoring

    array=("prometheus-prometheus-instance" )
  namespace=monitoring
  export STATUS_SUCCESS="Running"
  for i in "${array[@]}"
    do 
        echo ""
        echo "------------------------------------------------------------------------"
        echo "Check $i"
        while :
        do
            FIND=$i
            STATUS_CHECK=$(kubectl get pods -n $namespace | grep "$FIND" | awk '{print $3;}' | sed 's/"//g' | sed 's/,//g')
            echo "Status: $STATUS_CHECK"
            STATUS_VERIFICATION=$(echo "$STATUS_CHECK" | grep $STATUS_SUCCESS)
            if [ "$STATUS_VERIFICATION" = "$STATUS_SUCCESS" ]; then
                echo "$(date +'%F %H:%M:%S') Status: $FIND is Ready"
                echo "------------------------------------------------------------------------"
                break
            else
                echo "$(date +'%F %H:%M:%S') Status: $FIND($STATUS_CHECK)"
                echo "------------------------------------------------------------------------"
            fi
            sleep 3
        done
    done 
    
    echo "Press any key to move on"
    read input
}

function getPrometheusPortForward () {
   kubectl get service -n monitoring
   kubectl port-forward service/prometheus-instance -n monitoring 9090
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
    kubectl get pods -n fontendoperator-system

}

# **********************************************************************************
# Execution
# **********************************************************************************

echo "************************************"
echo " Install Prometheus Operator"
echo "************************************"
installPrometheusOperator

echo "************************************"
echo " Configure Prometheus Operator"
echo "************************************"
configurePrometheusOperator

echo "************************************"
echo " Build and upload custom operator container image"
echo "************************************"
# Uncomment the next line to use your own build
# ===============================
# buildAndUploadCustomOperator

echo "************************************"
echo " Deploy custom operator"
echo "************************************"
deployCustomOperator

echo "************************************"
echo " Port forward prometheus UI"
echo "************************************"
getPrometheusPortForward 

