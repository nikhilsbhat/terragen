# config


[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsbhat/terragen)](https://goreportcard.com/report/github.com/nikhilsbhat/terragen)  [![shields](https://img.shields.io/badge/license-mit-brightgreen](https://github.com/nikhilsbhat/terragen/blob/master/LICENSE)


A utility to generate the templates for building sutom [Terraform](https://www.terraform.io/) providers.

## Introduction

It is difficult to switch context of different kubernets clusters hosted in GCP projects.
If one has to connect cluster using gcloud, will end up runnig multiple gcloud commands and is painful task.

Yeah GCP has a option of cloud shell, where one can connect to the cluster hassle-free. Its little hard if we have to connect locally from our machines.

Config solves exactly the same thing, by letting one to switch the cluster in one command. At a stage it's interactive shell helps one in selection of the cluster they want to switch. As a bonus it also helps in activating service account and switching projects.

## Requires

As there are no prebuilt libraries, Terragen expects [GO](https://golang.org/dl/) installed on the machine to build one. This will help you on installing [GO](https://golang.org/doc/install)

## Installation

```golang
go get -u github.com/nikhilsbhat/terragen
go build
```
Use the executable just like any other go-cli application.

If incase few to use this in your piece of code import package in your code.
```golang
import (
    "github.com/nikhilsbhat/terragen"
)
```

### config commands

```bash
config [command] [flags]
```
Make sure appropriate command is used for the actions, to check the available commands and flags use `config --help`

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

### `config generate`

Credentials of GCP can be fed to `config` in two ways.
Either by passing path of credential file while invoking it or by setting environment variable `GOOGLE_APPLICATION_CREDENTIALS` just like how `gcloud` expects to be.

To switch to the cluster in the appropriate GCP you wish, `set` command helps in it.

```bash
config set -j /path/to/credential.json
    or
config set
```

You know which cluster to connect and don't want `config` to figure that out for you? then below command will help you with that

```bash
config set -c core-search-dev-cluster -r us-central1 -j /path/to/credential.json
```

**Note** `config set` without credentials file works only if the `GOOGLE_APPLICATION_CREDENTIALS` is set.

## Limitations

Right now this works only with kube clusters hosted in [GCP](https://cloud.google.com/), making it available accross other cloud will be more helpful.