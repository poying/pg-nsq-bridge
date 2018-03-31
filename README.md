# pg-nsq-bridge

Publish messages from Postgres to NSQ

### Usage

```bash
$ pg-nsq-bridge -h
NAME:
   pg-nsq-bridge

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.0.0

DESCRIPTION:
   Publish messages from Postgres to NSQ

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --pgchannel value   Postgres Channel Name (default: "pgnsqbridge")
   --nsqchannel value  NSQ Channel Name (default: "pgnsqbridge")
   --pgaddr value      Postgres Address (default: "host=127.0.0.1")
   --nsqaddr value     NSQ Address (default: "127.0.0.1:4150")
   --help, -h          show help
   --version, -v       print the version
```

### Install

```bash
$ git clone git@github.com:poying/pg-nsq-bridge.git
$ cd pg-nsq-bridge
$ make install
```