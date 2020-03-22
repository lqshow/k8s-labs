SHELL := /bin/bash

ifneq ($(NAMESPACE),)
	NAMESPACE ?= $(NAMESPACE)
else
	NAMESPACE ?= default
endif

SERVICE_ACCOUNT := $(SERVICE_ACCOUNT)
KUBE_APISERVER	:= $(KUBE_APISERVER)
CLUSTER_NAME	:= $(CLUSTER_NAME)
CONTEXT_NAME	:= $(CONTEXT_NAME)
ROLEBINDING_NAME:= $(ROLEBINDING_NAME) 	
ROLE_NAME		:= $(ROLE_NAME)
USER_ENTRY		:= $(SERVICE_ACCOUNT)-$(CLUSTER_NAME)
KUBECONFIG_FILE := k8s-$(SERVICE_ACCOUNT)-config

SECRET := $$(kubectl -n $(NAMESPACE) get secrets | grep ^$(SERVICE_ACCOUNT) | cut -f1 -d ' ')
USER_TOKEN := $$(kubectl -n $(NAMESPACE) get secret $(SECRET) -o go-template='{{.data.token}}' |base64 --decode)
CA_CRT_FILE := $(SERVICE_ACCOUNT)-ca.crt
TOKEN_FILE := $(SERVICE_ACCOUNT).token

# Create service account
create-sa:
	@kubectl -n $(NAMESPACE) create sa $(SERVICE_ACCOUNT)
	@echo "secret = $(SECRET)"

# Fetch the crt from the secret
create-ca-crt:
	@kubectl -n $(NAMESPACE) get secret $(SECRET) -o yaml | awk '/ca.crt:/{print $$2}' |base64 --decode > $(CA_CRT_FILE)
	@echo "----------->$(SERVICE_ACCOUNT)-ca.crt"
	@cat $(CA_CRT_FILE)
	@echo "<----------$(SERVICE_ACCOUNT)-ca.crt"

# Fetch the token from the secret
create-token:
	@kubectl -n $(NAMESPACE) get secret $(SECRET) -o go-template='{{.data.token}}' |base64 --decode > $(TOKEN_FILE)
	@echo "----------->$(SERVICE_ACCOUNT).token"
	@cat $(TOKEN_FILE)
	@echo "<----------$(SERVICE_ACCOUNT).token"

# Sets a cluster entry in kubeconfig
set-cluster:
	@kubectl config set-cluster $(CLUSTER_NAME) \
				--certificate-authority=./$(CA_CRT_FILE) \
				--embed-certs=true \
				--server=${KUBE_APISERVER} \
				--kubeconfig=$(KUBECONFIG_FILE)

# Sets a user entry in kubeconfig
set-credentials:
	@kubectl config set-credentials $(USER_ENTRY) \
			    --token=$(USER_TOKEN) \
			    --kubeconfig=$(KUBECONFIG_FILE)

# Sets a context entry in kubeconfig
set-context:
	@kubectl config set-context $(CONTEXT_NAME) \
			    --cluster=$(CLUSTER_NAME) \
			    --user=$(USER_ENTRY) \
			    --namespace=$(NAMESPACE) \
			    --kubeconfig=$(KUBECONFIG_FILE)

# Sets the current-context in a kubeconfig file
set-current-context:
	@kubectl config use-context $(CONTEXT_NAME) --kubeconfig=$(KUBECONFIG_FILE)

# Create a RoleBinding for a particular Role or ClusterRole
create-rolebinding:
	@kubectl -n $(NAMESPACE) create rolebinding $(ROLEBINDING_NAME) \
			    --role=$(ROLE_NAME) \
			    --serviceaccount=$(NAMESPACE):$(SERVICE_ACCOUNT)

# Create clusterrolebinding
create-clusterrolebinding:
	@kubectl -n $(NAMESPACE) create clusterrolebinding $(ROLEBINDING_NAME) \
			    --clusterrole=cluster-admin \
			    --serviceaccount=$(NAMESPACE):$(SERVICE_ACCOUNT)

create-serviceaccount: create-sa create-ca-crt create-token
	@echo  "The Service Account <$(SERVICE_ACCOUNT)> is created successfully"

setup-kubeconfig: set-cluster set-credentials set-context set-current-context
	@echo "Set up kube config <$(CLUSTER_NAME)> successfully"

.PHONY: create-serviceaccount \
		setup-kubeconfig \
		create-sa \
		create-ca-crt \
		create-token \
		set-cluster \
		set-credentials \
		set-context
