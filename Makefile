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
	@command -v docker >/dev/null 2>&1 || { echo >&2 "docker is not installed. Please install it first."; exit 1; }

.PHONY: create-cluster
create-cluster: check-requirements
	@if ! kind get clusters | grep -q "^$(CLUSTER_NAME)$$"; then \
		echo "Creating kind cluster $(CLUSTER_NAME)..."; \
		kind create cluster --name $(CLUSTER_NAME) --wait=60s --config cluster/config.yaml --kubeconfig $(KUBECONFIG_PATH); \
	else \
		echo "Cluster $(CLUSTER_NAME) already exists. Skipping creation."; \
	fi

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

install-%-operator:
	@echo "Installing $* operator..."
	cd $* && make && cd ..
	docker build -t $*controller -f $*/Dockerfile .
	kind load docker-image $*controller --name $(CLUSTER_NAME)
	kubectl apply -k $*/config/default

.PHONY: install-operators
install-operators: install-app-operator install-route-operator install-maintenance-operator

# Full deploy: create cluster, wait, apply kustomization
.PHONY: deploy
deploy: create-cluster apply-installation install-operators
