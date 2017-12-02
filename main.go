package main

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"github.com/9seconds/ddoctor/internal/checkers"
	"github.com/9seconds/ddoctor/internal/config"
)

var (
	app = kingpin.New(
		"ddoctor",
		"Docker doctor - checking health of your containers.")

	debug = app.Flag("debug", "Run in debug mode.").
		Short('d').
		Envar("DDOCTOR_DEBUG").
		Bool()
	oneShot = app.Flag("one-shot", "Do not run forever, execute only once").
		Short('o').
		Bool()

	configFile = app.Arg("config-path", "Path to the config").
			Required().
			File()
)

func init() {
	app.Version("0.0.1")
	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.WarnLevel)
}

func main() {
	app.Parse(os.Args[1:])

	if *debug {
		log.SetLevel(log.DebugLevel)
	}

	cf, err := config.ParseConfigFile(*configFile)
	if err != nil {
		log.Fatalf(err.Error())
	}

	log.WithFields(log.Fields{
		"periodicity": cf.Periodicity.Duration,
		"host":        cf.Host,
		"port":        cf.Port,
	}).Info("Config file")
	for _, v := range cf.Checks {
		log.WithFields(log.Fields{
			"type":    v.Type,
			"url":     v.URL.URL,
			"exec":    v.Exec,
			"timeout": v.Timeout.Duration,
		}).Info("Check")
	}

	channel := make(chan *checkers.CheckResult)
	checker, _ := checkers.NewShellChecker(1*time.Second, "ls -la")
	go checker.Run(context.Background(), channel)

	fmt.Println(<-channel)
}
