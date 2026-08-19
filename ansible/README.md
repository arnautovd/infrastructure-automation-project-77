# Ansible

`terraform-secrets.yml` is an Ansible Vault encrypted file. It contains the names of secrets used by Terraform, not plaintext credentials.

Create or edit it with:

```text
ansible-vault edit terraform-secrets.yml
```

Expected variables:

```yaml
vscale_token: "..."
terraform_backend_access_key: "..."
terraform_backend_secret_key: "..."
```

The Vault password must be supplied externally through `ANSIBLE_VAULT_PASSWORD_FILE` or the CI secret store.
