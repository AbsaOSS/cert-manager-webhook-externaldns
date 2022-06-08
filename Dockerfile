FROM alpine:3.14.0

RUN apk add --no-cache ca-certificates

COPY webhook /usr/local/bin/webhook

RUN apk add libcap && setcap 'cap_net_bind_service=+ep' /usr/local/bin/webhook

USER 1001

ENTRYPOINT ["webhook"]
