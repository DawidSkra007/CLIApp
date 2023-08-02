#!/bin/bash

echo "Stopping and removing containers..."
docker-compose down

echo "Cleaning up vault-data directories..."
rm -rf ./vault/*
rm -rf ./vault-data2/*

echo "Starting containers..."
docker-compose up --build -d

echo "Waiting for containers to start..."
sleep 10

echo "Initializing Keycloak..."
source ./keycloak_init.sh

echo "Initializing and unsealing Vault..."
./init_unseal_vault.sh

echo "Initializing and unsealing Vault..."
./init_unseal_vault2.sh
