// Code generated by gen-routes.
// DO NOT EDIT

package main

import (
	"flag"
	"path/filepath"
	"strings"

	"github.com/gophergala2016/source/internal"
)

var _ = filepath.Separator

const (
	appName = "source"
)

var (
	port = flag.String("port", "", "Port number to listen the application. (default \"8888\")")
	sock = flag.String("sock", "", "Path to a UNIX socket to listen the application.")
)

func main() {
	flag.Parse()
	if len(*port) == 0 {
		*port = ":8888"
	}

	if !strings.HasPrefix(*port, ":") {
		*port = ":" + *port
	}

	// Initialize App
	internal.Init(appName, *port)

	router := router()

	router.Static("/public", internal.JoinPath("./client"))

	router.LoadHTMLGlob(filepath.Clean(internal.JoinPath("./client")),
		filepath.Clean(internal.JoinPath("./client")))

	switch {
	case len(*sock) != 0:
		router.RunUnix(*sock)
	default:
		router.Run(*port)
	}
}
