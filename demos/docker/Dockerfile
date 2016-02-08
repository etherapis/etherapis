FROM alpine:3.3

MAINTAINER Péter Szilágyi <peterke@gmail.com>

# Pull in demo playground dependencies
RUN \
  apk add --update python py-pip libstdc++             && \
  apk add --update python-dev build-base linux-headers && \
  pip install circus                                   && \
  apk del python-dev build-base linux-headers          && \
  rm -rf /var/cache/apk/*

# Configure the Go runtime
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
ENV GO15VENDOREXPERIMENT=1

# Build all the demo services and clean up build environment
RUN \
  apk add --update go git gcc musl-dev                  && \
  go get github.com/etherapis/etherapis/demos/geolookup && \
  go get github.com/etherapis/etherapis/demos/streamer  && \
  go get github.com/etherapis/etherapis/etherapis       && \
  apk del go git gcc musl-dev                           && \
  rm -rf /var/cache/apk/* $GOPATH/src $GOPATH/pkg

# Initialize any demo services having external dependencies
RUN \
  mkdir -p /demos/geolookup /demos/streamer /etherapis && \
  streamer --init --root /demos/streamer/data

EXPOSE 80
EXPOSE 81
EXPOSE 8080
EXPOSE 8081

ADD circus.ini circus.ini
ENTRYPOINT ["circusd", "circus.ini"]
