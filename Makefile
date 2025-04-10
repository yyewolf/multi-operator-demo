# Name of the kind cluster
CLUSTER_NAME ?= multi-operator-demo

# Kustomization directory
KUSTOMIZE_DIR := cluster/prerequisites

# Kubeconfig path
KUBECONFIG_PATH := cluster/kubeconfig

# Default target
.PHONY: all
all: deploy

# Check if required tools are installed
.PHONY: check-requirements
check-requirements:
	@command -v kind >/dev/null 2>&1 || { echo >&2 "kind is not installed. Please install it first."; exit 1; }
	@command -v kubectl >/dev/null 2>&1 || { echo >&2 "kubectl is not installed. Please install it first."; exit 1; }
	@command -v helm >/dev/null 2>&1 || { echo >&2 "helm is not installed. Please install it first."; exit 1; }

# Create a kind cluster
.PHONY: create-cluster
create-cluster: check-requirements
	kind create cluster --name $(CLUSTER_NAME) --wait=60s --config cluster/config.yaml --kubeconfig $(KUBECONFIG_PATH)

# Delete the kind cluster
.PHONY: delete-cluster
delete-cluster: check-requirements
	kind delete cluster --name $(CLUSTER_NAME)

# Apply the installation to the kind cluster
.PHONY: apply-installation
apply-installation: check-requirements
	./cluster/prerequisites/install.sh
	@echo "Installation applied successfully."
	@echo "You can now access the cluster using the kubeconfig at $(KUBECONFIG_PATH)."

# Full deploy: create cluster, wait, apply kustomization
.PHONY: deploy
deploy: create-cluster apply-kustomization
