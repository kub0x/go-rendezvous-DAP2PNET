FROM golang:1.15-buster AS base

ENV CC=gcc
ENV GO111MODULE=on

RUN apt-get update && \
  apt-get --no-install-recommends install -y gcc \
    libc6-dev-armel-cross && \
  rm -rf /var/lib/apt/lists/*

ARG PKG=rendezvous 
ARG CA
ARG CERTIFICATE
ARG PRIVATE_KEY

RUN mkdir -p /go/src \
 && mkdir -p /go/bin \
 && mkdir -p /go/pkg

ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH

RUN mkdir -p $GOPATH/src/app
ADD . $GOPATH/src/app

WORKDIR $GOPATH/src/app

RUN mkdir -p /go/bin/certs/

RUN echo "$CA" | tee /go/bin/certs/ca.pem
RUN echo "$CERTIFICATE" | tee /go/bin/certs/$PKG.dap2p.net.pem
RUN echo "$PRIVATE_KEY" | tee /go/bin/certs/$PKG.dap2p.net.key

COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY *.go ./

RUN ls /usr/bin/gcc*

ENV CGO_ENABLED=1
RUN CC=gcc go build -o /go/bin/$PKG -race

WORKDIR /go/bin

# FROM scratch

# WORKDIR /

# EXPOSE 6667
# USER 1001

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder --chown=1001 /go/bin/$PKG ./$PKG
# COPY --from=builder --chown=1001 /go/bin/certs/ ./certs

ENTRYPOINT ["./rendezvous"]
