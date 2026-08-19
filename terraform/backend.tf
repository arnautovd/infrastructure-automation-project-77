terraform {
  backend "remote" {
    workspaces {
      name = "infrastructure-automation-project-77"
    }
  }
}
