package main

import (
	"flag"
	"fmt"

	"github.com/elijk/nibo/chi"
)

var flagHost *string = flag.String("h", "localhost", "host to listen to")
var flagPort *int = flag.Int("p", 3300, "port to listen on")
var flagConfig *string = flag.String("c", "", "path to configuration file (json)")

func main() {
	flag.Parse()

	// Load configuration from flags (and file, if provided).
	conf, err := loadConfig()
	if err != nil {
		panic(err)
	}

	// Setup server dependencies and additional configuration.
	opts := chi.ServerOpts{
		Configuration: chi.Configuration{
			Host: conf.Host,
			Port: conf.Port,
		},
	}

	// Create and start server.
	srv, err := chi.NewServer(opts)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Nibo server is now listening on %s:%d...\n", conf.Host, conf.Port)
	srv.Start()
}
