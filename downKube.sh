#!/bin/bash

echo "Stopping and removing Kubernetes resources..."
kubectl delete -f k8s-manifests/ --ignore-not-found

echo "Cleaning up vault-data directories..."
rm -rf ./vault/*
rm -rf ./vault-data2/*

killall kubectl