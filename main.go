package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/poying/pg-nsq-bridge/bridge"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

const version = "0.0.1"

func main() {
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt)

	app := cli.NewApp()
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "channel",
			Value: "pgchannel:nsqchannel",
			Usage: "Postgres & NSQ Channel Name",
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
		cli.BoolFlag{
			Name:  "raw",
			Usage: "Output JSON format",
		},
	}
	app.Name = "pg-nsq-bridge"
	app.Usage = ""
	app.Description = "Publish messages from Postgres to NSQ"
	app.Version = version
	app.Action = func(c *cli.Context) error {
		if c.Bool("raw") {
			log.SetFormatter(&log.JSONFormatter{})
		}
		options := bridge.Options{
			Channel:         c.String("channel"),
			PostgresAddress: c.String("pgaddr"),
			NSQAddress:      c.String("nsqaddr"),
			OnError: func(err error) {
				log.Error(err)
			},
			OnProcess: func(postgresChannel, nsqChannel string, payload bridge.Payload) {
				fields := log.Fields{
					"postgresChannel": postgresChannel,
					"nsqChannel":      nsqChannel,
					"payload":         payload,
				}
				log.WithFields(fields).Info("Receive message")
			},
		}
		b, err := bridge.New(options)
		if err != nil {
			return err
		}
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			<-exitCh
			cancel()
		}()
		b.Start(ctx)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
