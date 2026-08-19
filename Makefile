SHELL := /bin/sh

TF_DIR := terraform
TF := terraform -chdir=$(TF_DIR)

.PHONY: init fmt validate plan apply destroy check

init:
	$(TF) init \
		-backend-config="endpoint=$(TF_STATE_ENDPOINT)" \
		-backend-config="access_key=$(TF_STATE_ACCESS_KEY)" \
		-backend-config="secret_key=$(TF_STATE_SECRET_KEY)" \
		-backend-config="bucket=$(TF_STATE_BUCKET)"

fmt:
	$(TF) fmt -recursive

validate: init
	$(TF) validate

plan: validate
	$(TF) plan -out=tfplan

apply: plan
	$(TF) apply tfplan

destroy: validate
	$(TF) destroy

check: fmt validate
