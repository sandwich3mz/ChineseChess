FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS builder

LABEL maintainer="sandwich3mz"

# 在容器根目录创建 src 目录
WORKDIR /src

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache g++

COPY ./go.mod .
COPY ./go.sum .

ENV GOPROXY="https://goproxy.cn"

RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLE=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -ldflags "-s -w" -o http_server ./main.go

FROM alpine:latest

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
# RUN apk add --no-cache tzdata

WORKDIR /app

COPY ./config.yaml /app/chess/config.yaml
COPY ./resource/res.csv /app/chess/resource/res.csv

COPY --from=builder /src/http_server /app/

EXPOSE 8083

ENTRYPOINT ["./http_server"]