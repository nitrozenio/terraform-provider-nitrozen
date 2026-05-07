---
page_title: "nitrozen_project Resource - nitrozen"
description: |-
  Manages a Nitrozen changelog project.
---

# nitrozen_project (Resource)

Manages a changelog project on [Nitrozen.io](https://nitrozen.io).

## Example Usage

```terraform
resource "nitrozen_project" "example" {
  name        = "My App"
  description = "Product updates and release notes"
}
```

## Argument Reference

- `name` (Required) — The name of the project.
- `description` (Optional) — An optional description for the project.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` — The unique integer ID of the project.

## Import

Existing projects can be imported using the project ID:

```
terraform import nitrozen_project.example 123
```
