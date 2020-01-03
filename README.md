# terragen


[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsbhat/terragen)](https://goreportcard.com/report/github.com/nikhilsbhat/terragen)  [![shields](https://img.shields.io/badge/license-mit-brightgreen)](https://github.com/nikhilsbhat/terragen/blob/master/LICENSE) [![shields](https://godoc.org/github.com/nikhilsbhat/terragen?status.svg)](https://godoc.org/github.com/nikhilsbhat/terragen)


A utility to generate the templates for building custom [Terraform](https://www.terraform.io/) providers.

## Introduction

Terraform is one of the best software available to automate infrastructure and no procrastination in accepting that.

How about extending Terraform? this would be great when we are solving a complexity that are project specific. And Terraform has the same feature in the form of [custom-provider](https://www.terraform.io/docs/extend/how-terraform-works.html). How to create one such provider? they have well [documented](https://www.terraform.io/docs/extend/writing-custom-providers.html) steps for it.

How about creating initial/basic things with single command? terragen addresses the same. It generates the templates which eases the development of the custom-provider

## Requires

* Since there are no prebuilt libraries of Terragen available, it expected that [go](https://golang.org/dl/) to be pre installed on the machine to build one. Installing go can be found [here](https://golang.org/doc/install).
* Understanding of how to build [custom-provider](https://www.terraform.io/docs/extend/writing-custom-providers.html) for terraform.

## Installation

```golang
go get -u github.com/nikhilsbhat/terragen
go build
```
Use the executable just like any other go-cli application.

Found some of the codes useful for you? then start using it by importing the package in your line of codes.
```golang
import (
    "github.com/nikhilsbhat/terragen"
)
```

### terragen commands

```bash
terragen [command] [flags]
```
Make sure appropriate command is used for the actions, to check the available commands and flags use `terragen --help`

```bash
Terragen helps user to create custom terraform provider by generating templates for it.

Usage:
  terragen [command] [flags]

Available Commands:
  generate    command to generate the initial components for terraform provider
  help        Help about any command
  version     command to fetch the version of terragen installed

Flags:
  -h, --help          help for terragen
  -n, --name string   name of the provider that needs templates
  -p, --path string   path where the templates has to be generated

Use "terragen [command] --help" for more information about a command."
```

### `terragen generate`

The command `generate` helps in generating the provider templates which can be enhanced.

```bash
terragen generate
```

Pass the required falgs to generate the templates in the required location and with the proper name. If not, it defaults to working directory with the previder name `demo` just like above command.

```bash
terragen generate -n demo -p ~/terraform/demo/
```

**Note** `terragen generate` just creates the templates not the whole provider.

## Limitations

Currently it just creates basic templates, more features to be added to make it generate fully working provider.