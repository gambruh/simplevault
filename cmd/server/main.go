package main

import (
	"log"
)

func main() {
	if err := newRootCmd(runServer).Execute(); err != nil {
		log.Fatal(err)
	}
}
