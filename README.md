# Infrastructure

Инфраструктурный код проекта для Vscale.

## Структура

- `terraform/` — конфигурация Terraform и удалённого S3 backend.
- `ansible/` — зашифрованные Ansible Vault переменные и локальные инструкции.

Секреты не хранятся в открытом виде. Перед запуском необходимо заменить значения-заглушки в Vault-файле и передать пароль Vault внешним способом.

## Terraform

Установите Terraform и выполните команды из директории `terraform`:

```text
terraform init \
  -backend-config="endpoint=$env:TF_STATE_ENDPOINT" \
  -backend-config="access_key=$env:TF_STATE_ACCESS_KEY" \
  -backend-config="secret_key=$env:TF_STATE_SECRET_KEY" \
  -backend-config="bucket=$env:TF_STATE_BUCKET"
terraform plan -var="vscale_token=$env:VSCALE_API_TOKEN"
```

Для PowerShell сначала задайте `$env:VSCALE_API_TOKEN`, а затем используйте `-var="vscale_token=$env:VSCALE_API_TOKEN"`. Токен из сообщения не добавляется в файлы.

Backend-бакет создаётся заранее в S3-совместимом объектном хранилище. Значение `key` задаётся в `terraform/backend.tf` и определяет путь к state-файлу. Учётные данные backend передаются только через `-backend-config` или переменные окружения.

## Ansible Vault

Зашифрованный файл находится в `ansible/terraform-secrets.yml`. Пароль Vault не коммитится:

```text
ansible-vault edit ansible/terraform-secrets.yml
ansible-vault view ansible/terraform-secrets.yml
```

В CI пароль передавайте через `ANSIBLE_VAULT_PASSWORD_FILE` или секрет CI. После расшифровки не коммитьте файл и не создавайте из него `terraform.tfvars` в репозитории.
