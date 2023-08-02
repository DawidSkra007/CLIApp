#!/bin/bash

echo "Stopping and removing Kubernetes resources..."
kubectl delete -f k8s-manifests/ --ignore-not-found

echo "Cleaning up vault-data directories..."
rm -rf ./vault/*
rm -rf ./vault-data2/*

echo "Deploying Kubernetes resources..."
kubectl apply -f k8s-manifests/

echo "Waiting for Keycloak pod to become ready..."
kubectl rollout status deployment/keycloak --timeout=300s

echo "Setting up proxy to the Kubernetes API server..."
kubectl proxy --port=8080 &
vault proxy --port=8400 &
PROXY_PID=$!
sleep 1
# Ensure the proxy process is terminated when the script exits
trap "echo 'Terminating Kubernetes API server proxy...'; kill $PROXY_PID" EXIT
# Wait for Keycloak to become available
echo "Waiting for Keycloak to become available..."

echo "Initializing Keycloak..."
source ./keycloak_init.sh

echo "Initializing and unsealing Vault..."
./init_unseal_vault.sh

echo "Initializing and unsealing Vault..."
./init_unseal_vault2.sh

echo "Terminating Kubernetes API server proxy..."
killall kubectl