TF_DIR := terraform
TF := terraform -chdir=$(TF_DIR)
PROVIDER_DIR := provider
PROVIDER_MIRROR := provider-mirror/registry.terraform.io/arnautovd/vscale/0.1.0/windows_amd64
PROVIDER_NAME := terraform-provider-vscale_v0.1.0_x5.exe
export TF_CLI_CONFIG_FILE := $(CURDIR)/.terraformrc

.PHONY: provider-build init fmt validate plan apply destroy check

provider-build:
	powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path '$(PROVIDER_MIRROR)' | Out-Null"
	go -C $(PROVIDER_DIR) build -o ../$(PROVIDER_MIRROR)/$(PROVIDER_NAME) .

init: provider-build
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
