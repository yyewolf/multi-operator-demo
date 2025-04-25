#!/bin/bash

set -euo pipefail
LOG_FILE="deploy.log"
exec > >(tee -a "$LOG_FILE") 2>&1

echo "[$(date)] Starting deployment script"

# Change to script directory
cd "$(dirname "$0")"
echo "KUBECONFIG: $KUBECONFIG"

# Apply Gateway API CRDs
echo "[$(date)] Applying Gateway API CRDs..."
kubectl apply -k gateway-api/config/crd || {
  echo "[$(date)] Failed to apply Gateway API CRDs"
  exit 1
}

# Install or upgrade envoy-gateway
echo "[$(date)] Installing or upgrading envoy-gateway..."
if helm status eg -n envoy-gateway-system > /dev/null 2>&1; then
  helm upgrade eg oci://docker.io/envoyproxy/gateway-helm \
    --version v0.0.0-latest -n envoy-gateway-system -f envoy-values.yaml
else
  helm install eg oci://docker.io/envoyproxy/gateway-helm \
    --version v0.0.0-latest -n envoy-gateway-system --create-namespace -f envoy-values.yaml
fi

# Wait for envoy-gateway to be available
echo "[$(date)] Waiting for envoy-gateway deployment to be available..."
kubectl wait --timeout=5m -n envoy-gateway-system deployment/envoy-gateway --for=condition=Available || {
  echo "[$(date)] envoy-gateway deployment did not become available in time"
  exit 1
}

# Apply gateway resources
echo "[$(date)] Applying gateway resources..."
kubectl apply -k gateway || {
  echo "[$(date)] Failed to apply gateway resources"
  exit 1
}

# Add or update Helm repo
echo "[$(date)] Adding/updating jetstack Helm repo..."
helm repo add jetstack https://charts.jetstack.io --force-update

# Install or upgrade cert-manager
echo "[$(date)] Installing or upgrading cert-manager..."
if helm status cert-manager -n cert-manager > /dev/null 2>&1; then
  helm upgrade cert-manager jetstack/cert-manager \
    --namespace cert-manager \
    --version v1.17.0 \
    --set crds.enabled=true
else
  helm install cert-manager jetstack/cert-manager \
    --namespace cert-manager \
    --create-namespace \
    --version v1.17.0 \
    --set crds.enabled=true
fi

# Change to app directory
echo "[$(date)] Moving to app directory..."
cd ../../app/docker || {
  echo "[$(date)] App directory not found"
  exit 1
}

# Build Docker image
echo "[$(date)] Building Docker image..."
docker build -t workload . || {
  echo "[$(date)] Docker build failed"
  exit 1
}

# Load Docker image into Kind
echo "[$(date)] Loading Docker image into Kind..."
kind load docker-image workload --name=multi-operator-demo || {
  echo "[$(date)] Failed to load Docker image into Kind"
  exit 1
}

echo "[$(date)] Deployment completed successfully"
