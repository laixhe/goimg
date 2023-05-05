FROM golang:alpine as builder

ENV GOPROXY https://proxy.golang.com.cn,direct
ENV CGO_ENABLED 0

WORKDIR /goimg
COPY . /goimg

RUN go mod download

# Build
RUN go build -o goimg

FROM alpine

ENV TZ Asia/Shanghai

WORKDIR /goimg

COPY --from=builder /goimg/goimg .
COPY --from=builder /goimg/config.yaml .

EXPOSE 8080

ENTRYPOINT ["./goimg"]

