---
layout: default
title: Resource
parent: Create
nav_order: 2
---

## Create Resource

Subcommand `resource` of `terragen create` generates scaffolds for a specified resource of a selected provider.

`terragen create resource --help` would help in listing all available flags of `terragen create resource`.

```shell
This will help user to generate scaffolds for resource of chosen provider.

Usage:
  terragen create resource [args] [flags]

Flags:
  -h, --help              help for resource
      --provider string   name of the provider for which resource/datasource to be scaffolded (default "demo")

Global Flags:
      --dry-run       dry-run the process of provider scaffold creation
  -p, --path string   path where the templates has to be generated (default ".")
```

### Usage

* More resources to an existing project can be registered by running `terragen create resource [resource-name] --provider [provider-name]`.
* Referencing to example chosen for the provider, the command would look like `terragen create resource hashicups_coffee_order --provider hashicups`. 
* With the addition of new resource the metadata of project should now look like:

```yaml
version: 1.0.0
repo-group: github.com/nikhilsbhat
project-module: github.com/nikhilsbhat/terraform-provider-hashicups
provider: hashicups
provider-path: /Users/sample/my-opensource/terraform-provider-test
resources:
  - hashicups_order
  - hashicups_coffee_order
data-sources:
  - hashicups_coffees
  - hashicups_ingredients
importers:
  - ""
```

### Configuration

| Flags                | Type     | Description                                                                                      | Defaults |
|:--------------------:|:--------:|:-------------------------------------------------------------------------------------------------|:--------:|
| `provider`{: .fs-3 } | `string` | name of previously scaffolded terraform `provider` to which the `resource` to be registered with.|  demo    |