package main

import (
  "fs2consul/internal/consulclient"
  "fs2consul/internal/fsread"
  "log"
  "os"
  "strings"
)

func main() {
  if len(os.Args) < 3 {
    log.Printf(`
      Usage:
        %[1]s fs_dir consul_prefix
      Example:
        %[1]s fs/consul/kv/dir/ /services/`, os.Args[0])
    os.Exit(1)
  }
  consulAddr, exist := os.LookupEnv("CONSUL_ADDR")
  consulToken, existToken := os.LookupEnv("CONSUL_TOKEN")
  if !(exist && existToken) {
    log.Fatal("Environment variables CONSUL_ADDR or CONSUL_TOKEN must be defined")
  }
  args := os.Args[1:]
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
  //log.Print("FOUND KVs: ")
  //s, err := json.MarshalIndent(fs.KV, "", "  ")
  //if err != nil {
  //    log.Printf("ERROR: %s\n", err)
  //}
  //log.Print(string(s))
  //log.Println(strings.Repeat("=", 100))
  consul := consulclient.NewConsulClient(consulAddr, consulToken, consulPrefix)
  //consul.ListConsulKVs()
  log.Println(strings.Repeat("=", 100))
  err := consul.SyncKV(fs.KV)
  if err != nil {
    log.Panic("PANIC: fs and KV not sync!!!")
  }
}