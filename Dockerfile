FROM alpine:3.6

RUN apk add --no-cache --update curl \
  dumb-init \
  bash \
  grep \
  sed \
  jq \
  ca-certificates \
  openssl && \
  rm -rf /var/cache/apk/*

ADD tldrfeed /bin/tldrfeed
RUN chmod +x /bin/tldrfeed
EXPOSE 8080

ENTRYPOINT "/bin/tldrfeed"