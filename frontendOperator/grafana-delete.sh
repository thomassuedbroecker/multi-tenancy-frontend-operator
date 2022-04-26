#!/bin/bash

kubectl delete -f grafana-instance.yaml -n grafana-operator
kubectl delete -f grafana-operator-setup.yaml -n grafana-operator