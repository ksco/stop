FROM golang:latest AS builder
RUN go env -w GOPROXY=https://goproxy.cn,direct
WORKDIR /backend
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/slack .

FROM alpine:latest

RUN apk add tzdata && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    apk del tzdata

COPY --from=builder /go/bin/slack /slack

WORKDIR /

ENTRYPOINT ["/slack"]