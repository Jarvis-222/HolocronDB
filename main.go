package main

import (
	"HolocronDB/config"
	"HolocronDB/server"
	"flag"
	"log"
)

func setupFlags() {

	flag.StringVar(&config.Host, "host", "0.0.0.0", "host for the HolocronDB server")
	flag.IntVar(&config.Port, "port", 7379, "port for the HolocronDB server")
	flag.Parse()

}

func main() {

	setupFlags()
	log.Printf("This is master Obi-Wan Kenobi, you are on the planet HolocronDB, ready to serve your requests on %s:%d", config.Host, config.Port)
	log.Println("Hello there !!")
	server.RunSyncTCPServer()
}
