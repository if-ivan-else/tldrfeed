FROM alpine:3.7

RUN apk add --no-cache ca-certificates && \
  rm -rf /var/cache/apk/*

ADD tldrfeed /bin/tldrfeed
RUN chmod +x /bin/tldrfeed
EXPOSE 8080
ENV DB_URL "mongo:27017/db"
CMD ["/bin/tldrfeed", "server"]