#!/bin/bash
# **************** Global variables
# TDB

# **********************************************************************************
# Functions
# **********************************************************************************


# Example Operator
function deleteFrontendOperator () {
    kubectl delete -f config/samples/multitenancy_v1alpha1_tenancyfrontend.yaml -n default
    make undeploy
    echo "Press any key to move on"
    read input
}


function deletePrometheusInstanceAndConfiguration () {
    # Prometheus Configuration
    #kubectl delete -f prom-serviceaccount.yaml -n monitoring
    #kubectl delete -f prom-clusterrole.yaml -n monitoring
    #kubectl delete -f prom-clusterrolebinding.yaml -n monitoring
    kubectl delete -f prom-instance.yaml -n monitoring
    kubectl delete -f prom-loadbalancer.yaml -n monitoring
    kubectl delete namespace monitoring
    echo "Press any key to move on"
    read input
}

# Prometheus Operator
function deletePrometheusOperator () {
    # kubectl delete -f https://operatorhub.io/install/prometheus.yaml
    kubectl delete -f ./prom-bundle-monitoring.yaml
    echo "Press any key to move on"
    read input
}

function deleteCRDsPrometheusOperator () {
    kubectl delete customresourcedefinition alertmanagerconfigs.monitoring.coreos.com
    kubectl delete customresourcedefinition podmonitors.monitoring.coreos.com
    kubectl delete customresourcedefinition servicemonitors.monitoring.coreos.com
    kubectl delete customresourcedefinition thanosrulers.monitoring.coreos.com
    kubectl delete customresourcedefinition prometheusrules.monitoring.coreos.com
    kubectl delete customresourcedefinition prometheuses.monitoring.coreos.com
    kubectl delete customresourcedefinition probes.monitoring.coreos.com
    kubectl delete customresourcedefinition alertmanagers.monitoring.coreos.com
    kubectl delete customresourcedefinition podmonitors.monitoring.coreos.com
    echo "Press any key to move on"
    read input
}


function deleteOLM () {
    # Delete OLM
    operator-sdk olm uninstall latest
    echo "Press any key to move on"
    read input
}

# **********************************************************************************
# Execution
# **********************************************************************************

deleteFrontendOperator
deletePrometheusInstanceAndConfiguration
deletePrometheusOperator
deleteCRDsPrometheusOperator
deleteOLM