# Terraform Provider for Nitrozen

Manage [Nitrozen.io](https://nitrozen.io) changelog projects and entries with Terraform.

[![Registry](https://img.shields.io/badge/terraform-registry-blue?logo=terraform)](https://registry.terraform.io/providers/nitrozenio/nitrozen)
[![GitHub](https://img.shields.io/badge/github-nitrozenio%2Fterraform--provider--nitrozen-blue?logo=github)](https://github.com/nitrozenio/terraform-provider-nitrozen)

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.23 (to build from source)

## Usage

```hcl
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
```

Get your API token from the [API Tokens](https://nitrozen.io/api-tokens) page.

## Resources

### `nitrozen_project`

Manages a changelog project.

```hcl
resource "nitrozen_project" "my_app" {
  name        = "My App"
  description = "Changelog for My App"
}
```

**Attributes:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `name` | string | yes | Project name |
| `description` | string | no | Project description |
| `id` | number | computed | Project ID |

---

### `nitrozen_entry`

Manages a changelog entry within a project.

```hcl
resource "nitrozen_entry" "v1_release" {
  project_id   = nitrozen_project.my_app.id
  title        = "v1.0 — Initial Release"
  content      = "We shipped v1.0! Here's what's included."
  category     = "new"
  is_published = true
}
```

**Attributes:**

| Name | Type | Required | Description |
|------|------|----------|-------------|
| `project_id` | number | yes | ID of the parent project (forces replacement) |
| `title` | string | yes | Entry title |
| `content` | string | yes | Entry body (Markdown supported) |
| `category` | string | yes | One of: `new`, `improvement`, `fix`, `announcement` |
| `is_published` | bool | no | Whether the entry is publicly visible. Default: `false` |
| `id` | number | computed | Entry ID |

## Complete Example

```hcl
resource "nitrozen_project" "api" {
  name        = "API"
  description = "API changelog"
}

resource "nitrozen_entry" "rate_limits" {
  project_id   = nitrozen_project.api.id
  title        = "Rate limit headers added"
  content      = "All API responses now include `X-RateLimit-Remaining` headers."
  category     = "improvement"
  is_published = true
}

resource "nitrozen_entry" "auth_fix" {
  project_id   = nitrozen_project.api.id
  title        = "Fixed token expiry bug"
  content      = "OAuth tokens were expiring 1 hour early. This is now fixed."
  category     = "fix"
  is_published = true
}
```

## Building from Source

```sh
git clone https://github.com/nitrozenio/terraform-provider-nitrozen
cd terraform-provider-nitrozen
go build -o terraform-provider-nitrozen
```

## License

MIT
