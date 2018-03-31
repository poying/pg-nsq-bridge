package main

import (
	"os"
	"os/signal"

	"github.com/poying/pg-nsq-bridge/bridge"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const version = "0.0.0"

func main() {
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt)

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "pgchannel",
			Value: "pgnsqbridge",
			Usage: "Postgres Channel Name",
		},
		cli.StringFlag{
			Name:  "nsqchannel",
			Value: "pgnsqbridge",
			Usage: "NSQ Channel Name",
		},
		cli.StringFlag{
			Name:  "pgaddr",
			Value: "host=127.0.0.1",
			Usage: "Postgres Address",
		},
		cli.StringFlag{
			Name:  "nsqaddr",
			Value: "127.0.0.1:4150",
			Usage: "NSQ Address",
		},
	}
	app.Name = "pg-nsq-bridge"
	app.Usage = ""
	app.Description = "Publish messages from Postgres to NSQ"
	app.Version = version
	app.Action = func(c *cli.Context) error {
		options := bridge.Options{
			PostgresChannel: c.String("pgchannel"),
			NSQChannel:      c.String("nsqchannel"),
			PostgresAddress: c.String("pgaddr"),
			NSQAddress:      c.String("nsqaddr"),
			OnError: func(err error) {
				log.Error(err)
			},
			OnProcess: func(payload bridge.Payload) {
				log.Info("Receive notifcation: " + payload.Data)
			},
		}
		b, err := bridge.New(options)
		if err != nil {
			return err
		}
		go b.Start()
		defer b.Stop()
		<-exitCh
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
