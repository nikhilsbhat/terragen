---
layout: default
title: Edit
nav_order: 3
has_children: true
---

## Edit

Command `edit` of terragen helps to modify terragen created provider, datasource and resource of terraform.

## terragen edit

`terragen edit --help` would help in listing all available sub commands and flags of `terragen edit`.

```shell
This will help user to edit the scaffolds generated for terraform provider and other components of them.

Usage:
  terragen edit [command] [flags]

Available Commands:
  datasource  Command to edit already generated scaffolds of a datasource
  provider    Command to edit already generated scaffolds of a provider
  resource    Command to edit already created scaffolds generated scaffolds of resource

Flags:
  -h, --help   help for edit

Global Flags:
      --dry-run       dry-run the process of provider scaffold creation
  -p, --path string   path where the templates has to be generated (default ".")


Use "terragen edit [command] --help" for more information about a command."
```

For more information editing these components navigate to their independent documents.