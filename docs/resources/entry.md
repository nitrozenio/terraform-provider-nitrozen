---
page_title: "nitrozen_entry Resource - nitrozen"
description: |-
  Manages a changelog entry within a Nitrozen project.
---

# nitrozen_entry (Resource)

Manages a changelog entry within a [Nitrozen.io](https://nitrozen.io) project. Entries are the individual items that appear in your public changelog.

## Example Usage

```terraform
resource "nitrozen_project" "example" {
  name = "My App"
}

resource "nitrozen_entry" "release" {
  project_id   = nitrozen_project.example.id
  title        = "v2.0 — Dark mode"
  content      = "We added dark mode support across the entire app. Toggle it from your profile settings."
  category     = "new"
  is_published = true
}
```

### Draft entry

```terraform
resource "nitrozen_entry" "draft" {
  project_id = nitrozen_project.example.id
  title      = "Upcoming maintenance window"
  content    = "Scheduled maintenance on Saturday 2am–4am UTC."
  category   = "announcement"
  # is_published defaults to false
}
```

## Argument Reference

- `project_id` (Required) — The ID of the parent project. Changing this forces a new resource.
- `title` (Required) — The title of the changelog entry.
- `content` (Required) — The body content of the entry. Markdown is supported.
- `category` (Required) — The category of the entry. Must be one of:
  - `new` — New features
  - `improvement` — Improvements to existing features
  - `fix` — Bug fixes
  - `announcement` — General announcements
- `is_published` (Optional) — Whether the entry is publicly visible. Defaults to `false`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

- `id` — The unique integer ID of the entry.

## Import

Existing entries can be imported using the project ID and entry ID separated by a slash:

```
terraform import nitrozen_entry.example 42/7
```
