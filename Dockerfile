FROM golang:1.26-alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0
ENV GOPROXY=https://goproxy.cn,direct

# tzdata 时区 + ca-certificates 根证书（运行时 TLS 验证必须）
RUN apk update --no-cache && apk add --no-cache tzdata ca-certificates

WORKDIR /build

ADD go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /app/main ./cmd/server/

# ---- 运行时镜像 ----
FROM alpine:3.19

# 直接在运行时镜像安装证书，比从 builder 复制更可靠
RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata

ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/main .

RUN mkdir -p /app/logs

EXPOSE 8099

CMD ["./main"]
