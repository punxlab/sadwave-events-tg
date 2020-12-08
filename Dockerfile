FROM alpine

EXPOSE 8080

COPY bin/sadwave-events-tg .
COPY config.json .

CMD ["/sadwave-events-tg"]