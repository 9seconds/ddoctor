package main

import (
	"fmt"
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	app = kingpin.New("ddoctor", "Docker doctor - checking health of your containers.")

	debug      = app.Flag("debug", "Run in debug mode.").Short('d').Envar("DDOCTOR_DEBUG").Bool()
	configPath = app.Flag("config-path", "Path to the config").Envar("DDOCTOR_CONFIG_PATH").Short('c').ExistingFile()
	oneShot    = app.Flag("one-shot", "Do not run forever, execute only once").Short('o').Bool()
)

func init() {
	app.Version("0.0.1")
}

func main() {
	app.Parse(os.Args[1:])
	fmt.Println(*debug)
}
