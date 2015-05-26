package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/sreis/kuproxy/haproxy"
	"github.com/sreis/kuproxy/keystore"
)

func main() {

	var master *string = flag.String("master", "http://172.17.8.101:2379", "Etcd master connection url.")
	flag.Parse()

	if err := haproxy.Start(); err != nil {
		log.Fatal("Error launching haproxy. '", err, "'.")
		return
	}

	haproxy.ShowStat()

	// Handle Ctrl-C
	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)

		<-sigchan

		if err := haproxy.Stop(); err != nil {
			log.Fatal("Error stopping haproxy. '", err, "'.")
		}

		os.Exit(0)
	}()

	log.Println("Watching and waiting for pods to come online.")
	log.Println("Press Ctrl-C to exit...")
	keystore.Watch(*master)
}
