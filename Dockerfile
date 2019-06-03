FROM golang:1.11.6

ADD ./webhook/webhook /webhook
RUN chmod +x /webhook
ENTRYPOINT ["/webhook"]
