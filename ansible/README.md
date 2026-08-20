# Ansible

Все Ansible-файлы находятся в этой директории. Основной playbook — `playbook.yml`.

## Подготовка

Установите зависимости из `requirements.yml`:

```text
ansible-galaxy role install -r requirements.yml -p roles
ansible-galaxy collection install -r requirements.yml -p collections
```

Inventory `inventory/hosts.yml` генерируется из Terraform outputs командой `make ansible-inventory`. Файл не коммитится.

## Секреты

`terraform-secrets.yml` — зашифрованный Ansible Vault-файл. Реальные токены и пароли не хранятся в репозитории.

Создайте отдельный пароль Vault вне репозитория и редактируйте файл:

```text
ansible-vault edit terraform-secrets.yml
```

Переменные окружения для SSH (`ANSIBLE_REMOTE_USER`, `ANSIBLE_PRIVATE_KEY_FILE`) и Docker image передаются снаружи.

## Теги

```text
ansible-playbook -i inventory/hosts.yml playbook.yml --tags prepare
ansible-playbook -i inventory/hosts.yml playbook.yml --tags deploy
ansible-playbook -i inventory/hosts.yml playbook.yml --tags verify
```
