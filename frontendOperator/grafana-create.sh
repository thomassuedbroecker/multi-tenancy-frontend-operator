#!/bin/bash

kubectl create -f grafana-instance.yaml -n grafana-operator
kubectl create -f grafana-operator-setup.yaml -n grafana-operator
