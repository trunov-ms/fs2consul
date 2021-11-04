package main

import (
	"fmt"
	"fs2consul/internal/consulclient"
	"fs2consul/internal/fsread"
	"log"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 4 {
		fmt.Printf(`
	  Requirements
	  Environment variables:
		CONSUL_ADDR - your Consul address. Ex.: http://consul.local:8500
		CONSUL_HTTP_TOKEN - token with rw permission on prefix path
	  Usage
		%[1]s <get-diff|apply> <fs_dir> <consul_prefix>
		get-diff show difference between Consul KV and directiry
		apply sync data
	  Example
		%[1]s get-diff fs/consul/kv/dir/ /services/%s`, os.Args[0], "\n")
		os.Exit(1)
	}
	consulAddr, exist := os.LookupEnv("CONSUL_ADDR")
	consulToken, existToken := os.LookupEnv("CONSUL_HTTP_TOKEN")
	if !(exist && existToken) {
		log.Fatal("Environment variables CONSUL_ADDR or CONSUL_HTTP_TOKEN must be defined")
	}
	command := os.Args[1]
	args := os.Args[2:]
	basedir := args[0]
	consulPrefix := args[1]
	if !strings.HasSuffix(consulPrefix, "/") {
		consulPrefix = consulPrefix + "/"
	}
	if !strings.HasSuffix(basedir, "/") {
		basedir = basedir + "/"
	}
	basedir = strings.TrimPrefix(basedir, "./")
	fs := fsread.NewFSRead(basedir)
	consul := consulclient.NewConsulClient(consulAddr, consulToken, consulPrefix)
	if strings.Compare(command, "apply") == 0 {
		err := consul.SyncKV(fs.KV, true)
		if err != nil {
			log.Print(err.Error())
			log.Panic("PANIC: fs and KV not sync!!!")
		}
	} else if strings.Compare(command, "get-diff") == 0 {
		consul.SyncKV(fs.KV, false)
	}

}
