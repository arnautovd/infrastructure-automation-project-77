TF_DIR := terraform
TF := terraform -chdir=$(TF_DIR)
PROVIDER_DIR := provider
PROVIDER_MIRROR := provider-mirror/registry.terraform.io/arnautovd/vscale/0.1.6/windows_amd64
PROVIDER_NAME := terraform-provider-vscale_v0.1.6_x5.exe
export TF_CLI_CONFIG_FILE := $(CURDIR)/.terraformrc

.PHONY: provider-build init fmt validate plan apply destroy check

provider-build:
	powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path '$(PROVIDER_MIRROR)' | Out-Null"
	go -C $(PROVIDER_DIR) build -o ../$(PROVIDER_MIRROR)/$(PROVIDER_NAME) .

init: provider-build
	$(TF) init -upgrade \
		-backend-config="organization=$(HCP_TERRAFORM_ORGANIZATION)"

fmt:
	$(TF) fmt -recursive

validate: init
	$(TF) validate

plan: validate
	$(TF) plan

apply: validate
	$(TF) apply

destroy: validate
	$(TF) destroy

check: fmt validate
