FROM alpine:3.5

COPY pg-nsq-bridge /usr/bin/pg-nsq-bridge

CMD [ "pg-nsq-bridge", "-h" ]