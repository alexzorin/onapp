OnApp API for Golang and CLI
========================

[![Build Status](https://secure.travis-ci.org/alexzorin/onapp.png?branch=master)](http://travis-ci.org/alexzorin/onapp)

This is to be a library for accessing the OnApp API with a Go interface.

In addition, there is to be a command line client that will mimic the feature set of the OnApp dashboard. Install via `go install github.com/alexzorin/onapp/onapp`

API
-------

Documentation for the API can be found at [godoc.org](http://godoc.org/github.com/alexzorin/onapp)

CLI
-------

You can install the `onapp` CLI command into your `$GOPATH/bin` via:

`go get github.com/alexzorin/onapp` + `go install github.com/alexzorin/onapp/onapp`

or download a binary release from here, if available.

Get started with `onapp config`, find usage via `onapp help` and `onapp help [command]`.

### Commands

* `config` - Configure the tool to connect to a particular OnApp Dashboard Server (saves by default in `~/.onapp`)
* `test`: Test the config
* `help`: Help text for all commands and subcommands
* `vm`: Management of virtual machines
    - `list <id>`: List virtual machines and their current status in a table
    - `start <id>`: Start a virtual machine
    - `stop <id>`: Stop a virtual machine
    - `reboot <id>`: Reboot a virtual machine
    - `ssh <id>`: Launches `ssh` at the VM's first IP address and provides you with the root password
    - `stat <id>`: SSH's into the machine (no password prompt) and runs `vmstat 1 10`, which it relays to `stdout`
    - `tx <id> [num_to_list]`: List of recent transactions on that VM

Where `<id>` is mentioned, you may either provide exact #ID, exact Label or Hostname, or the CLI will attempt to gues which VM you mean via Levenshtein distance. Inexact matches will prompt confirmation.