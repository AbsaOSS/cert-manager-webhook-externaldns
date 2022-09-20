FROM alpine:3.15.0

RUN apk add --no-cache ca-certificates
WORKDIR /usr/local/bin/

COPY webhook /usr/local/bin/webhook
# COPY apiserver.local.config/certificates/apiserver.crt /tls/tls.crt
# COPY apiserver.local.config/certificates/apiserver.key /tls/key.crt

RUN apk add libcap && setcap 'cap_net_bind_service=+ep' /usr/local/bin/webhook

USER 1001

ENTRYPOINT ["webhook"]
