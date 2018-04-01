SRC := main.go $(shell find bridge -type f)
DESTDIR := /usr/local/bin
VERSION := "0.0.1"

build: pg-nsq-bridge

clean:
	@rm -rf pg-nsq-bridge

install: build
	@mv pg-nsq-bridge $(DESTDIR)

image:
	@make clean
	@env GOOS=linux GOARCH=arm make build
	@docker build -t poying/pg-nsq-bridge:$(VERSION) .
	@docker push poying/pg-nsq-bridge:$(VERSION) 

pg-nsq-bridge: $(SRC)
	@go build -o $@ main.go