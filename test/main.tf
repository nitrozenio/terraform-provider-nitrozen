terraform {
  required_providers {
    nitrozen = {
      source = "nitrozenio/nitrozen"
    }
  }
}

# Set NITROZEN_TOKEN env var or pass token directly (do not commit tokens)
provider "nitrozen" {
  token = var.nitrozen_token
}

variable "nitrozen_token" {
  description = "Nitrozen API token"
  type        = string
  sensitive   = true
}

resource "nitrozen_project" "example" {
  name        = "My Changelog"
  description = "Product updates and release notes"
}

resource "nitrozen_entry" "example" {
  project_id   = nitrozen_project.example.id
  title        = "Initial release"
  content      = "We launched! Here's what's included."
  category     = "new"
  is_published = true
}
