package main

import (
	"flag"
)

var flagRunAddr string
var reportInterval int
var pollInterval int

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.IntVar(&reportInterval, "r", 10, "frequency of sending metrics to the server")
	flag.IntVar(&pollInterval, "p", 2, "frequency of gathering metrics")
	flag.Parse()
}
