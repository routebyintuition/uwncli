[![Go Report Card](https://goreportcard.com/badge/github.com/routebyintuition/uwncli)](https://goreportcard.com/badge/github.com/routebyintuition/uwncli) [![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0) [![Build Status](https://travis-ci.com/routebyintuition/uwncli.svg?branch=main)](https://travis-ci.com/routebyintuition/uwncli) [![codecov](https://codecov.io/gh/routebyintuition/uwncli/branch/main/graph/badge.svg?token=149HXZ5XZY)](https://codecov.io/gh/routebyintuition/uwncli)

# uwncli
Unikum und Wunderbar Nutanix CLI

UWNCLI is a Nutanix command line utility build to replicate and build upon the existing NCLI (Nutanix CLI) capabilities. The existing NCLI available for download via your Nutanix Prism Central/Element interface runs on Java where this application is build using Go. This is a work in progress with limited capabilities but consistently expanding as features are requested.

---

- [uwncli](#uwncli)
  - [Install](#install)
  - [Build](#build)
  - [Configuration](#configuration)
  - [Capabilities](#capabilities)
  - [Examples](#examples)

---

## Install
Download the [latest release](https://github.com/routebyintuition/uwncli/releases) specific to your operating system. There are packages for Linux, Mac, and Windows. 

```sh
tar -zxvf uwncli_<version number>_<operating system>_x86_64.tar.gz
./uwncli
```

## Build
The main branch is currently the development branch which will change soon. For now, it is recommended to use the [uwncli releases](https://github.com/routebyintuition/uwncli/releases)
```sh
go get -u github.com/routebyintuition/uwncli
```

## Configuration
uwncli can be configured to use stored credentials. These credentials are saved by default in ~/.nutanix/. This is a ".nutanix" folder under your user home directory. To configure saved credentials, we have an example. When entering text, all input is treated as sensitive and not echoed to the screen.

```sh
#> uwncli configure
Profile name [default]:
prism central username [admin]:
prism central password []:
prism central address [10.0.0.1:9440]:
invalid input length. less than 6 characters
saved profile to:  /Users/<username>/.nutanix/test.credential
```

Authentication and configuration can also be done via exported environmental variables as well as command line flags as shown below.

Exported command line variables:
```sh
export NUTANIX_PC_ADDRESS="10.0.0.10:9440"
export NUTANIX_PC_USER="username"
export NUTANIX_PC_PASS="password"
uwncli vm list
```

Command line flags:
```sh
uwncli --pcaddress "10.0.0.10:9440" --username <username> --password <password> vm list
```

Skipping certificate verification can be useful for non-production environments or new deployments where a valid vertificate has not yet been configured. This can be done as in the example below:

```sh
./uwncli --skip-cert-verify vm list
```

## Capabilities

- configure
- profile
  - list
  - delete
  - create
- vm
  - list
  - get
  - disklist
  - update-memory
  - update-power
- disk
  - list
  - list-vdisk
- cluster
  - list
- image
  - list
  - create
- subnet
  - list
- karbon
  - cluster
    - list

## Examples

