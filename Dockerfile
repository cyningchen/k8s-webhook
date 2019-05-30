FROM golang:1.11.6

ADD ./webhook /
RUN chmod +x /webhook
ENTRYPOINT ["/webhook"]
