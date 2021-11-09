## Introduction

If you work with Consul KV you probably want to change it in VCS.
One of the ways to "represent" Consul KV storage in VCS files is a directory tree, with keys are the names of files and values are these contents.
This util synchronize the Consul KV storage and FS directory tree.

All changes apply in single transaction with [Consul Txn API](https://www.consul.io/api/txn)

## Usage

`fs2consul <get-diff|apply> <./path/to/dir/> <consul/kv/prefix>`

## Example

```bash
export CONSUL_HTTP_ADDR=http://localhost:8500
export CONSUL_HTTP_TOKEN=aaaaaaaa-bbbb-cccc-dddd-eeeeeeeeeeee
fs2consul apply ./consul-kv.git/ /services/
```
