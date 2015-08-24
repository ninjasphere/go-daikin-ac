package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ninjasphere/go-daikin-ac"
)

func main() {

	found, err := daikin.Discover(time.Second * 10)
	if err != nil {
		panic(err)
	}

	for ac := range found {
		spew.Dump("Found AC", ac)
	}
}
