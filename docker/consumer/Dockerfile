FROM golang:1.11 as builder

RUN mkdir -p /go/src/github.com/sergeiten/hugecsv
WORKDIR /go/src/github.com/sergeiten/hugecsv
COPY . .

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o app cmd/consumer/consumer.go

FROM alpine:latest
RUN set -ex && apk add --no-cache ca-certificates
WORKDIR /
COPY --from=builder /go/src/github.com/sergeiten/hugecsv/app .
CMD ./app
