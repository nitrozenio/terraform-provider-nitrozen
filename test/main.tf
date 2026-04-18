terraform {
  required_providers {
    nitrozen = {
      source = "girishadhithya/nitrozen"
    }
  }
}

provider "nitrozen" {
  token = "18|FwMwS1mqD0YL7KoHP78R5WNVmpLnqgjwVyf0u60x9564b830"
}

resource "nitrozen_project" "test" {
  name        = "new Name"
  description = "new Description"
}