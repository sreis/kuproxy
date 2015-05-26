package main

import (
	"flag"
	"log"

	"github.com/sreis/kuproxy/haproxy"
	"github.com/sreis/kuproxy/keystore"
)

//
func main() {

	var master *string = flag.String("master", "http://172.17.8.101:2379", "Etcd master connection url.")
	flag.Parse()

	if err := haproxy.Start(); err != nil {
		log.Fatal(err)
		return
	}

	keystore.Watch(*master)

	haproxy.Stop()
}
