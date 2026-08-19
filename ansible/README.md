# Ansible

`terraform-secrets.yml` is an Ansible Vault encrypted file. It contains the names of secrets used by Terraform, not plaintext credentials.

Create or edit it with:

```text
ansible-vault edit terraform-secrets.yml
```

Expected variables:

```yaml
vscale_token: "..."
hcp_terraform_token: "..."
```

The Vault password must be supplied externally through `ANSIBLE_VAULT_PASSWORD_FILE` or the CI secret store.
