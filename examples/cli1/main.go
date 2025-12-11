package main

import (
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

var dochost string = ""

func buildAppInstance() (appInst *cli.App) {
	appInst = &cli.App{
		Name:     "deploy_docs",
		Version:  "0.0.2",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Mike Pennington",
				Email: "mike@pennington.net",
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "dochost",
				Value:       "127.0.0.1",
				Usage:       "FQDN or IPv4 of the documentation host",
				Destination: &dochost,
			},
		},
		Action: func(cCtx *cli.Context) (err error) {
			log.Println("Starting cli.Context action")
			if cCtx.NArg() == 0 {
				log.Fatal("No CLI arguments detected!")
			}
			log.Printf("args: %+v", cCtx.Args())
			return nil
		},
	}
	return appInst
}

func main() {
	app := buildAppInstance()
	log.Printf("dochost before Run: %q", dochost)
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
	log.Printf("dochost after Run: %q", dochost)
}
