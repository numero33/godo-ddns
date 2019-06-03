FROM golang:alpine
ADD . /go/src/godo-ddns/
WORKDIR /go/src/godo-ddns/

RUN apk --update --no-cache add git \
&& rm -rf /var/cache/apk/* \
&& go get ./... \
&& go build \
&& chmod +x /go/bin/godo-ddns \
&& rm -rf /go/src
ENTRYPOINT ["/go/bin/godo-ddns"]