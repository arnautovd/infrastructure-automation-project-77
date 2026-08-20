$ErrorActionPreference = "Stop"

$output = terraform -chdir=terraform output -json web_server_addresses | ConvertFrom-Json
if ($output -is [array]) {
    $addresses = @($output)
} else {
    $addresses = @($output.value)
}
if ($addresses.Count -ne 2 -or $addresses -contains $null -or $addresses -contains "") {
    throw "Terraform must output two non-empty web server addresses."
}

$lines = @(
    "---",
    "all:",
    "  children:",
    "    web:",
    "      hosts:"
)

for ($index = 0; $index -lt $addresses.Count; $index++) {
    $name = "web-{0:D2}" -f ($index + 1)
    $lines += "        ${name}:"
    $lines += "          ansible_host: $($addresses[$index])"
}

Set-Content -LiteralPath "ansible/inventory/hosts.yml" -Value $lines -Encoding utf8
Write-Output "Generated ansible/inventory/hosts.yml for $($addresses.Count) web servers."
