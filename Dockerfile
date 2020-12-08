FROM alpine

EXPOSE 80

COPY bin/sadwave-events-tg /bin
COPY config.json .

CMD ["/bin/sadwave-events-tg"]