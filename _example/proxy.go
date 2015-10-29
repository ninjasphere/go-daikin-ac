package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/ninjasphere/go-daikin-ac/proxy"
	"github.com/ninjasphere/go-ninja/config"
)

var udpBroadcast = config.String("192.168.12.255", "daikin-proxy.broadcast") + ":30050"
var iface = config.String("eth2", "daikin-proxy.interface")

func main() {

	err := proxy.Start(udpBroadcast, iface)
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)
}
