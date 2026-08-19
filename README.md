# Infrastructure

Инфраструктурный код проекта для Vscale.

## Структура

- `terraform/` — конфигурация Terraform и удалённого HCP Terraform backend.
- `ansible/` — зашифрованные Ansible Vault переменные и локальные инструкции.
- `Makefile` — инициализация, проверка, планирование и управление инфраструктурой.
- `provider/` — собственный Vscale provider на Go и Terraform Plugin Framework.

Секреты не хранятся в открытом виде. Перед запуском необходимо заменить значения-заглушки в Vault-файле и передать пароль Vault внешним способом.

## Terraform

Установите Terraform, GNU Make и Go 1.23+.

Создайте workspace `infrastructure-automation-project-77` в HCP Terraform и задайте во внешних переменных:

```text
terraform init \
  -backend-config="organization=$env:HCP_TERRAFORM_ORGANIZATION"
terraform plan -var="vscale_token=$env:VSCALE_API_TOKEN"
```

Для PowerShell сначала задайте `$env:VSCALE_API_TOKEN`, а затем используйте `-var="vscale_token=$env:VSCALE_API_TOKEN"`. Токен из сообщения не добавляется в файлы.

Для HCP Terraform задайте также `$env:TF_TOKEN_app_terraform_io`. Этот токен не добавляется в репозиторий. Terraform state будет храниться в workspace HCP Terraform.

HCP token должен иметь доступ к workspace как минимум с правом `Plan`; для `make apply` требуется право записи/применения. Если workspace использует наборы permissions, выдайте токену соответствующую роль.

## Создаваемая инфраструктура

Terraform создаёт ровно две VM Vscale с именами `web-01` и `web-02`. Образ, тариф, локация и ID заранее зарегистрированных SSH-ключей задаются переменными Terraform.

В выбранном варианте Vscale не создаются Managed Load Balancer и Managed Database. Балансировщик и база данных должны быть подключены отдельно, если они понадобятся приложению.

Часто используемые команды:

```text
make init
make check
make plan
make apply
make destroy
```

`make init` сначала собирает custom provider из `provider/`, затем инициализирует Terraform через локальный filesystem mirror и HCP Terraform backend. Provider не скачивается из Terraform Registry.

Проверка без создания ресурсов:

```text
make check
make plan
```

`make plan` обращается к удалённому backend и формирует локальный файл `terraform/tfplan`. Создание ресурсов выполняется только командой `make apply`.

Для PowerShell можно передать токен без записи в файлы:

```powershell
$env:TF_VAR_vscale_token = $env:VSCALE_API_TOKEN
make plan
```

SSH-ключи должны быть заранее добавлены в Vscale. Их ID передаются, например, через `TF_VAR_vscale_ssh_key_ids='["123"]'`.

## Ansible Vault

Зашифрованный файл находится в `ansible/terraform-secrets.yml`. Пароль Vault не коммитится:

```text
ansible-vault edit ansible/terraform-secrets.yml
ansible-vault view ansible/terraform-secrets.yml
```

В CI пароль передавайте через `ANSIBLE_VAULT_PASSWORD_FILE` или секрет CI. После расшифровки не коммитьте файл и не создавайте из него `terraform.tfvars` в репозитории.
