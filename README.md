# config


[![Go Report Card](https://goreportcard.com/badge/github.com/nikhilsbhat/terragen)](https://goreportcard.com/report/github.com/nikhilsbhat/terragen)  [![shields](https://img.shields.io/badge/license-mit-brightgreen](https://github.com/nikhilsbhat/terragen/blob/master/LICENSE)


A utility which helps to switch between multiple [Kubernetes](https://kubernetes.io/) clusters of [GKE](https://cloud.google.com/kubernetes-engine/).

## Introduction

It is difficult to switch context of different kubernets clusters hosted in GCP projects.
If one has to connect cluster using gcloud, will end up runnig multiple gcloud commands and is painful task.

Yeah GCP has a option of cloud shell, where one can connect to the cluster hassle-free. Its little hard if we have to connect locally from our machines.

Config solves exactly the same thing, by letting one to switch the cluster in one command. At a stage it's interactive shell helps one in selection of the cluster they want to switch. As a bonus it also helps in activating service account and switching projects.

## Requires

This isn't a standalone tool it still depends on few things like [`gcloud`](https://cloud.google.com/sdk/gcloud/). But makes life lot easier handling it.
* [`gcloud`](https://cloud.google.com/sdk/install) version 253.0.0 or higher (tested)

## Installation

```golang
go get -u github.com/nikhilsbhat/config
go build
```
Use the executable just like any other go-cli application.

If incase few to use this in your piece of code import package in your code.
```golang
import (
    "github.com/nikhilsbhat/config"
)
```

### config commands

```bash
config [command] [flags]
```
Make sure appropriate command is used for the actions, to check the available commands and flags use `config --help`

```bash
This will help user to deal with gcloud and kube config activity.


Usage:
  config [command] [flags]

Available Commands:
  help        Help about any command
  set         command to set the config
  version     command to fetch the version of config installed

Flags:
  -c, --cluster-name string   name of the cluster which needs to be connected to
  -h, --help                  help for config
  -j, --json string           path to gcp auth json file
  -r, --region strings        region where your cluster resides
  -v, --version string        version of the cluster (default "1")

Use "config [command] --help" for more information about a command."
```

### `config set`

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