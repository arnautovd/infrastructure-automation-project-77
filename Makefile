TF_DIR := terraform
TF := terraform -chdir=$(TF_DIR)
PROVIDER_DIR := provider
PROVIDER_MIRROR := provider-mirror/registry.terraform.io/arnautovd/vscale/0.1.6/windows_amd64
PROVIDER_NAME := terraform-provider-vscale_v0.1.6_x5.exe
export TF_CLI_CONFIG_FILE := $(CURDIR)/.terraformrc
export ANSIBLE_CONFIG := $(CURDIR)/ansible/ansible.cfg
export ANSIBLE_ROLES_PATH := $(CURDIR)/ansible/roles
export ANSIBLE_COLLECTIONS_PATH := $(CURDIR)/ansible/collections

.PHONY: provider-build init fmt validate plan apply destroy check ansible-requirements ansible-inventory deploy deploy-prepare deploy-app deploy-verify

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

ansible-requirements:
	ansible-galaxy role install -r ansible/requirements.yml -p ansible/roles
	ansible-galaxy collection install -r ansible/requirements.yml -p ansible/collections

ansible-inventory:
	powershell -NoProfile -ExecutionPolicy Bypass -File ansible/scripts/generate_inventory.ps1

deploy-prepare: ansible-requirements ansible-inventory
	ansible-playbook -i ansible/inventory/hosts.yml ansible/playbook.yml --tags prepare --user "$(ANSIBLE_REMOTE_USER)" --private-key "$(ANSIBLE_PRIVATE_KEY_FILE)"

deploy-app: ansible-requirements ansible-inventory
	ansible-playbook -i ansible/inventory/hosts.yml ansible/playbook.yml --tags deploy --user "$(ANSIBLE_REMOTE_USER)" --private-key "$(ANSIBLE_PRIVATE_KEY_FILE)"

deploy-verify: ansible-inventory
	ansible-playbook -i ansible/inventory/hosts.yml ansible/playbook.yml --tags verify --user "$(ANSIBLE_REMOTE_USER)" --private-key "$(ANSIBLE_PRIVATE_KEY_FILE)"

deploy: deploy-prepare deploy-app deploy-verify
