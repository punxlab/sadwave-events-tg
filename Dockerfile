FROM alpine

EXPOSE 80

COPY bin/sadwave-events-tg .
COPY config.json .

CMD ["/sadwave-events-tg"]