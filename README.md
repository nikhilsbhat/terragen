# terragen

[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsbhat/terragen)](https://goreportcard.com/report/github.com/nikhilsbhat/terragen)
[![shields](https://img.shields.io/badge/license-mit-brightgreen)](https://github.com/nikhilsbhat/terragen/blob/master/LICENSE)
[![shields](https://godoc.org/github.com/nikhilsbhat/terragen?status.svg)](https://godoc.org/github.com/nikhilsbhat/terragen)
[![shields](https://img.shields.io/github/v/tag/nikhilsbhat/terragen.svg)](https://github.com/nikhilsbhat/terragen/tags)
[![shields](https://img.shields.io/github/downloads/nikhilsbhat/terragen/total.svg)](https://github.com/nikhilsbhat/terragen/releases)

A utility to ease the development of [terraform](https://www.terraform.io/) custom provider by generating scaffolds for provider and their components.

## Introduction

Terraform is one of the best software available to automate infrastructure and no procrastination in accepting this fact.<br><br>
How about extending Terraform?. It would be great when we try solving complexity that is project-specific, and Terraform offers the same feature in the form of a [custom provider](https://www.terraform.io/docs/extend/how-terraform-works.html). How to create one such provider? They have well [documented](https://www.terraform.io/docs/extend/writing-custom-providers.html) steps for it.<br><br>
How about accelerating the development of such custom-provider with the scaffolds? `Terragen` helps here. It generates scaffolds for `provider`, `resources`, and `data_sources` that eases the development of the custom provider.<br><br>
Supports the addition of new scaffolds of `data_source` and `resource` for a specific Terraform `provider` as and when required.  
## Requirements

* [Go](https://golang.org/dl/) 1.19 or above . Installing go can be found [here](https://golang.org/doc/install).
* Basic understanding of terraform provider and [building](https://www.terraform.io/docs/extend/writing-custom-providers.html) them.

## Installation

* Recommend installing released versions. Release binaries are available on the [releases](https://github.com/nikhilsbhat/terragen/releases) page.
* Can always build it locally by running `go build` against cloned repo.

## Features supported by the `Terragen` at the moment.

|  component   |    create  |     edit     |  delete  |
| :----------: | :--------: | :----------: | :------: |
| `provider`   | yes        | yes (beta)   | no       |
| `data_source`| yes        | yes (beta)   | no       |
| `resource`   | yes        | yes (beta)   | no       |
| `importer`   | no         | no           | no       |

## Documentation

* Document of `Terragen` on its usage is up [here](https://nikhilsbhat.github.io/terragen).
* Updated documentation on all available commands and flags can be found [here](https://github.com/nikhilsbhat/terragen/blob/master/cli_docs/doc/terragen.md).

## Installation

* Recommend installing released versions. Release binaries are available on the [releases](https://github.com/nikhilsbhat/terragen/releases) page and docker from [here](https://hub.docker.com/repository/docker/basnik/terragen).
* Can always build it locally by running `go build` against cloned repo.

#### Docker

Latest version of docker images are published to [ghcr.io](https://github.com/nikhilsbhat/terragen/pkgs/container/terragen), all available images can be found there. </br>

```bash
docker pull ghcr.io/nikhilsbhat/terragen:latest
docker pull ghcr.io/nikhilsbhat/terragen:<github-release-tag>
```

## TODO
* [ ] Edit feature to cover more aspect.
* [ ] Test cases.

**Note** `terragen generate` just creates the templates not the whole provider.