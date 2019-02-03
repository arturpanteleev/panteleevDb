package main

import (
	"flag"
	"panteleevDb/api"
)

func main() {
	bindAddr := flag.String("bind_addr", ":8080", "Set bind address")
	size := flag.Int("lru_size", 20, "Set lru size per driver")
	flag.Parse()
	a := api.New(*bindAddr, *size)
	a.Start()
	a.WaitStop()
}
