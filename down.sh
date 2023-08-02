#!/bin/bash

echo "Stopping and removing containers..."
docker-compose down

echo "Cleaning up vault-data directories..."
rm -rf ./vault/*
rm -rf ./vault-data2/*