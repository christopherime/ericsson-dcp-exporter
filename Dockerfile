FROM golang:1.19-alpine AS builder

ARG ARCH=amd64

ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PATH $GOPATH/bin:$GOROOT/bin:$PATH
ENV GO_VERSION 1.19
ENV GO111MODULE on
ENV CGO_ENABLED=0

# Build dependencies
WORKDIR /go/src/
COPY /src/ .
RUN apk update && apk add make git
RUN mkdir /go/src/build
RUN go build -o build/ericsson-dcp-exporter

# Second stage
FROM alpine:latest

COPY --from=builder /go/src/build/ericsson-dcp-exporter /usr/local/bin/ericsson-dcp-exporter
RUN mkdir /etc/ericsson-dcp-exporter
CMD ["/usr/local/bin/ericsson-dcp-exporter"]
EXPOSE 9742
