#!/bin/bash

echo "-> Git Clone"
echo ""
# git clone https://github.com/thomassuedbroecker/multi-tenancy-frontend-operator
# cd multi-tenancy-frontend-operator/frontendOperator
# git checkout "update-operator"
echo ""

echo ""
echo "-> Build contoller manager image"
echo ""

export VERSION=0.0.2
export REGISTRY='quay.io'
export ORG='tsuedbroecker'
export CONTROLLER_IMAGE='frontendcontroller:v4'

echo ""
echo "-> build"
echo "-> $REGISTRY/$ORG/$CONTROLLER_IMAGE"
echo ""

make generate
make manifests
# operator-sdk generate kustomize manifests
make docker-build IMG="$REGISTRY/$ORG/$CONTROLLER_IMAGE"

echo ""
echo "-> push"
echo ""
docker push "$REGISTRY/$ORG/$CONTROLLER_IMAGE"

echo ""
echo "-> create bundle"
make bundle IMG="$REGISTRY/$ORG/$CONTROLLER_IMAGE"
echo ""

echo ""
echo "-> build bundle"
export BUNDLE_IMAGE='bundlefrontendoperator:v4'
make bundle-build BUNDLE_IMG="$REGISTRY/$ORG/$BUNDLE_IMAGE"
echo ""

# export REGISTRY='quay.io'
# export ORG='tsuedbroecker'
# export BUNDLE_IMAGE='bundlefrontendoperator:v4'
# echo "Image: $REGISTRY/$ORG/$BUNDLE_IMAGE"
# operator-sdk run bundle "$REGISTRY/$ORG/$BUNDLE_IMAGE" -n operators

echo ""
echo "-> push bundle"
docker push "$REGISTRY/$ORG/$BUNDLE_IMAGE"
echo ""


echo ""
echo "-> build catalog"
export CATALOG_IMAGE=frontend-catalog
export CATALOG_TAG=v0.0.2
make catalog-build CATALOG_IMG="$REGISTRY/$ORG/$CATALOG_IMAGE:$CATALOG_TAG" BUNDLE_IMGS="$REGISTRY/$ORG/$BUNDLE_IMAGE"
echo ""

echo ""
echo "-> push catalog"
docker push "$REGISTRY/$ORG/$CATALOG_IMAGE:$CATALOG_TAG"
echo ""

echo ""
echo "-> kubernetes"
kubectl apply -f "./olm-configuration/catalogsource.yaml" -n operators
# kubectl delete -f "./olm-configuration/catalogsource.yaml" -n operators
kubectl apply -f "./olm-configuration/subscription.yaml" -n operators
# kubectl delete -f "./olm-configuration/subscription.yaml" -n operators
kubectl apply -f config/samples/multitenancy_v2alpha2_tenancyfrontend.yaml -n default
# kubectl delete -f config/samples/multitenancy_v2alpha2_tenancyfrontend.yaml -n default

#kubectl get catalogsource  frontend-operator-catalog-2 -n operators -oyaml
#kubectl get subscriptions frontendoperator-v0-0-2-sub -n operators -oyaml 