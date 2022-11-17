FROM alpine:latest

RUN apk --update --no-cache add ca-certificates && rm -rf /var/cache/apk/*

COPY godo-ddns /godo-ddns

RUN chmod +x /godo-ddns

ENTRYPOINT ["/godo-ddns"]