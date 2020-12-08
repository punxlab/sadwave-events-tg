FROM alpine

EXPOSE 80

COPY bin/ /bin
COPY config.json .

CMD ["/bin/sadwave-events-tg"]