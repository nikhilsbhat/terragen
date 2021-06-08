---
layout: default
title: Create
nav_order: 2
has_children: true
---

## Create

Command `create` of terragen helps in creation of provider, datasource and resource of terraform.

## terragen create

`terragen create --help` would help in listing all available sub commands and flags of `terragen create`.

```shell
This will help user to generate the initial components of terraform provider.

Usage:
  terragen create [command] [flags]

Available Commands:
  datasource  Command to generate scaffolds for datasource
  provider    Command to generate scaffolds for terraform provider
  resource    Command to generate scaffolds for resource

Flags:
  -h, --help   help for create

Global Flags:
      --dry-run       dry-run the process of provider scaffold creation
  -p, --path string   path where the templates has to be generated (default ".")


Use "terragen create [command] --help" for more information about a command."
```

For more information on how to create these components, navigate to their independent documents.

### Configuration

| Flags                | Description                                    | Defaults  |
|:--------------------:|:----------------------------------------------:|:---------:|
| `path`{: .fs-3 }     | path under which scaffolds to be generated     | cwd       |
| `dry-run`{: .fs-3 }  | this simulates the process by not creating one | false     |
