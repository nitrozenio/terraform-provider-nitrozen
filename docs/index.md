---
page_title: "Provider: Nitrozen"
description: |-
  Use the Nitrozen provider to manage changelog projects and entries on Nitrozen.io.
---

# Nitrozen Provider

The Nitrozen provider lets you manage [Nitrozen.io](https://nitrozen.io) changelog projects and entries using Terraform. Automate your release notes, keep changelogs in sync with your infrastructure, and manage entries as code.

## Example Usage

```terraform
terraform {
  required_providers {
    nitrozen = {
      source  = "nitrozenio/nitrozen"
      version = "~> 1.0"
    }
  }
}

provider "nitrozen" {
  token = var.nitrozen_token
}

variable "nitrozen_token" {
  description = "Nitrozen API token"
  type        = string
  sensitive   = true
}

resource "nitrozen_project" "my_app" {
  name        = "My App"
  description = "Changelog for My App"
}

resource "nitrozen_entry" "launch" {
  project_id   = nitrozen_project.my_app.id
  title        = "v1.0 — Initial Release"
  content      = "We shipped v1.0!"
  category     = "new"
  is_published = true
}
```

## Authentication

Create an API token from the [API Tokens](https://nitrozen.io/api-tokens) page in your Nitrozen account, then provide it to the provider:

```terraform
provider "nitrozen" {
  token = var.nitrozen_token
}
```

It is recommended to pass the token via a Terraform variable or environment variable rather than hardcoding it.

## Argument Reference

- `token` (Required, Sensitive) — Your Nitrozen API token.
