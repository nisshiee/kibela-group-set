FROM golang:1.12.4-alpine3.9 as build

RUN set -ex \
    \
    && apk update && apk upgrade \
    && apk add --no-cache git

RUN mkdir -p /go/src/github.com/nisshiee/kibela-group-set
WORKDIR /go/src/github.com/nisshiee/kibela-group-set
ENV GO111MODULE=on

COPY go.mod go.sum /go/src/github.com/nisshiee/kibela-group-set/
RUN go mod download

COPY . /go/src/github.com/nisshiee/kibela-group-set
RUN go build -o kibela-group-set *.go

FROM alpine:3.9

RUN set -ex \
    \
    && apk update && apk upgrade \
    && apk add --no-cache ca-certificates

ENV SSL_CERT_FILE /etc/ssl/certs/ca-certificates.crt

ENV KIBELA_TEAM ""
ENV KIBELA_TOKEN ""
ENV KIBELA_TARGET_FOLDER_ID ""
ENV KIBELA_TARGET_GROUP_ID ""

ENTRYPOINT ["/kibela-group-set"]
COPY --from=build /go/src/github.com/nisshiee/kibela-group-set/kibela-group-set /kibela-group-set
