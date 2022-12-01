FROM golang:alpine as builder

WORKDIR /app
COPY . /app

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct && \
    go mod tidy && \
    go build -o amazing

FROM alpine
WORKDIR /app
COPY --from=builder /app/amazing .
COPY config.yaml .
CMD ["/app/amazing"]