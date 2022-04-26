#!/bin/bash

kubectl delete -f grafana-instance.yaml -n grafana-operator
kubectl delete -f grafana-operator-setup.yaml -n grafana-operator

kubectl get crds | greb 'grafa'

kubectl delete customresourcedefinition grafanadashboards.integreatly.org
kubectl delete customresourcedefinition grafanadatasources.integreatly.org
kubectl delete customresourcedefinition grafananotificationchannels.integreatly.org
kubectl delete customresourcedefinition grafanas.integreatly.org 