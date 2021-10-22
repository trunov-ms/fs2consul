package consulclient

import (
	"encoding/json"
	fs "fs2consul/internal/fsread"
	"log"
	"strings"

	"github.com/hashicorp/consul/api"
)

var client *api.Client
var err error

type ConsulClient struct {
	addr   string
	token  string
	prefix string
}

type KVVerb struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
	Verb  string `json:"Verb"`
}

type KVTxn struct {
	KV KVVerb
}

func NewConsulClient(addr string, token string, prefix string) *ConsulClient {
	config := api.DefaultConfig()
	config.Address = addr
	config.Token = token
	config.Datacenter = "vcdev"
	client, err = api.NewClient(config)
	if err != nil {
		log.Fatalf("Fatal: %s\n", err.Error())
	}
	return &ConsulClient{
		addr:   addr,
		token:  token,
		prefix: prefix,
	}
}

// isKeyInArray check Consul key exist in fsread kv array
// key: consul key
// kv: fs KVPair fileName->Content
func isKeyInArray(key string, kv []fs.KVPair) bool {
	for _, v := range kv {
		if strings.Contains(key, v.Key) {
			return true
		}
	}
	return false
}

// isKVInArray check file exist in Consul KV
// k - fs KVPair fileName->Content
// consulKV - list all Consul KV with prefix
func isKVInArray(k fs.KVPair, consulKV []fs.KVPair) bool {
	for _, kv := range consulKV {
		if strings.Contains(kv.Key, k.Key) && k.Value == kv.Value {
			return true
		}
	}
	return false
}

// getKVTxnOps format fs.KVPair list to Consul KVTxnOps for transaction commit
// kv - fs KVPair fileName->Content
// prefix - start element for Consul KV-tree
func getKVTxnOps(kv []fs.KVPair, prefix string) api.KVTxnOps {
	var output []*api.KVTxnOp
	consulKV := getConsulKV(prefix)
	// compare Consul KV with fs
	for _, consulKv := range consulKV {
		// if key exist in Consul but not fs delete them
		if !isKeyInArray(consulKv.Key, kv) {
			var op api.KVTxnOp
			op.Key = consulKv.Key
			op.Value = []byte(consulKv.Value)
			op.Verb = api.KVDelete
			output = append(output, &op)
		}
	}
	// compare fs with Consul
	for _, kv := range kv {
		// if key or value in fs differ from Consul add or replace them
		if !isKVInArray(kv, consulKV) {
			var op api.KVTxnOp
			op.Key = prefix + kv.Key
			op.Value = []byte(kv.Value)
			op.Verb = api.KVSet
			output = append(output, &op)
		}
	}
	return output
}

// getConsulKV getting current KV-tree from Consul with prefix
func getConsulKV(prefix string) []fs.KVPair {
	kv := client.KV()
	var output []fs.KVPair
	data, meta, er := kv.List(strings.TrimSuffix(prefix, "/"), nil)
	if er != nil {
		log.Println(er)
		log.Println(meta)
		log.Println(data)
		log.Fatalf("FATAL: %s\n", er.Error())
	}
	if data == nil {
		log.Println("Empty data")
	} else {
		for _, kv := range data {
			k := fs.KVPair{
				Key:   kv.Key,
				Value: string(kv.Value),
			}
			output = append(output, k)
		}
	}
	return output
}

// SyncKV apply transaction in Consul KV
// https://www.consul.io/api-docs/txn
// kv - list of fs KVPair fileName->Content
func (cc *ConsulClient) SyncKV(kv []fs.KVPair) error {
	data := getKVTxnOps(kv, cc.prefix)
	if len(data) == 0 {
		log.Println("INFO: Nothing to commit. FS and KV are sync!!!")
		return nil
	}
	log.Println("INFO: These changes will be apply")
	s, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
	}
	log.Print(string(s))
	log.Println("TRANSACTION START")
	ckv := client.KV()
	ok, response, _, err := ckv.Txn(data, nil)
	if err != nil {
		log.Println("FATAL: " + err.Error())
		return err
	}
	log.Printf("STATE: %t\n", ok)
	if !ok {
		log.Fatal("ERROR: TRANSACTION ROLLED BACK")
		return err
	}
	log.Println("TRANSACTION COMMIT")
	log.Println("APPLIED KV KEYS:")
	s, err = json.MarshalIndent(response.Results, "", "  ")
	if err != nil {
		log.Fatalf("ERROR: %s\n", err.Error())
		return err
	}
	log.Print(string(s))
	for _, res := range response.Errors {
		log.Printf("%s\n", res.What)
	}
	return nil
}

func (cc *ConsulClient) ListConsulKVs() {
	kv := getConsulKV(cc.prefix)
	log.Println("CURRENT KV KEYS LIST")
	for _, k := range kv {
		log.Printf("%s: %s\n", k.Key, k.Value)
	}
	log.Println(strings.Repeat("=", 50))
}
