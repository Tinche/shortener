FROM golang:1.9-alpine3.7 AS builder

RUN apk add --no-cache curl git
RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.3.2/dep-linux-amd64 && chmod +x /usr/local/bin/dep

RUN mkdir -p /go/src/github.com
WORKDIR /go/src/github.com

COPY src/shortener/Gopkg.toml src/shortener/Gopkg.lock ./

RUN dep ensure -vendor-only

WORKDIR /go/src/app
COPY src/shortener/ .
RUN go build -o shortener *.go

FROM alpine:3.7

WORKDIR /root/
COPY --from=builder /go/src/app/shortener /root/shortener
ENTRYPOINT ["/root/shortener"]