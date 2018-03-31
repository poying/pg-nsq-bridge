SRC := main.go $(shell find bridge -type f)
DESTDIR := /usr/local/bin

build: pg-nsq-bridge

install: build
	@mv pg-nsq-bridge $(DESTDIR)

pg-nsq-bridge: $(SRC)
	@go build -o $@ main.go