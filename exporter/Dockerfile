FROM golang:alpine3.9 as builder
MAINTAINER Thomas Maier <contact@thomas-maier.net>
RUN apk add --no-cache git
COPY . /build
WORKDIR /build
RUN go build

FROM alpine:3.9
COPY --from=builder /build/exporter /exporter
CMD /exporter

