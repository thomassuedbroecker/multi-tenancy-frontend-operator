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

function deleteCustomOperator () {
   # delete custom operator and example application deployed by the custom operator
   cd ..
   kubectl delete -f config/samples/multitenancy_v1alpha1_tenancyfrontend.yaml -n default
   export REGISTRY='quay.io'
   export ORG='tsuedbroecker'
   export CONTROLLER_IMAGE='frontendcontroller-monitoring:v0.0.1'  
   make undeploy IMG="$REGISTRY/$ORG/$CONTROLLER_IMAGE"
   cd ./prometheus-example
   kubectl delete namespace frontendoperator-system
}

function deleteExampleAppAndPrometheusInstance () {

    echo "1-app-example-deployment.yaml"
    kubectl delete -f 1-app-example-deployment.yaml -n default
    echo "2-app-example-service.yaml"
    kubectl delete -f 2-app-service.yaml -n default
    echo "3-prom-service-monitor.yaml"
    kubectl delete -f 3-prom-service-monitor.yaml -n default
    echo "4-prom-service-account.yaml"
    kubectl delete -f 4-prom-service-account.yaml -n default
    echo "5-prom-cluster-role.yaml"
    kubectl delete -f 5-prom-cluster-role.yaml -n default
    echo "6-prom-cluster-role-binding.yaml"
    kubectl delete -f 6-prom-cluster-role-binding.yaml -n default
    echo "7-prom-instance-servicemonitor-selector-configured.yaml"
    kubectl delete -f 7-prom-instance-servicemonitor-selector-configured.yaml -n default
    echo "8-prom-expose-ui-loadbalancer.yaml"
    kubectl delete -f 8-prom-expose-ui-loadbalancer.yaml -n default
}

function deletePrometheusOperatorAndOLM () {

   # delete the example app from https://github.com/prometheus-operator/prometheus-operator/blob/main/Documentation/user-guides/getting-started.md
   kubectl delete namespace example-app

   # delete the prometheus operator
   kubectl delete -f https://operatorhub.io/install/prometheus.yaml 
   
   # delete remaining monitoring crds
   kubectl delete customresourcedefinition alertmanagerconfigs.monitoring.coreos.com
   kubectl delete customresourcedefinition podmonitors.monitoring.coreos.com
   kubectl delete customresourcedefinition servicemonitors.monitoring.coreos.com
   kubectl delete customresourcedefinition thanosrulers.monitoring.coreos.com
   kubectl delete customresourcedefinition prometheusrules.monitoring.coreos.com
   kubectl delete customresourcedefinition prometheuses.monitoring.coreos.com
   kubectl delete customresourcedefinition probes.monitoring.coreos.com
   kubectl delete customresourcedefinition alertmanagers.monitoring.coreos.com
   kubectl delete customresourcedefinition podmonitors.monitoring.coreos.com

   # uninstall olm
   operator-sdk olm uninstall latest

   # delete custerroles not automated
   # ================================
   # RESULT=$(kubectl get clusterrole | grep "prom" | awk '{print $1;}'| head -n 1)
   # echo "Result: $RESULT"
   # kubectl delete clusterrole $RESULT
   # kubectl get clusterrole | grep "prom"
}


# **********************************************************************************
# Execution
# **********************************************************************************

echo "************************************"
echo " connectToIBMCluster"
echo "************************************"

connectToIBMCluster

echo "************************************"
echo " deleteCustomOperator"
echo "************************************"

deleteCustomOperator

echo "************************************"
echo " deleteExampleAppAndPrometheusInstance"
echo "************************************"

deleteExampleAppAndPrometheusInstance

echo "************************************"
echo " deletePrometheusOperatorAndOLM"
echo "************************************"

deletePrometheusOperatorAndOLM