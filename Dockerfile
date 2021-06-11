FROM golang:1.15.13-buster as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN go mod download

COPY . /workspace

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

FROM alpine:3.13
WORKDIR /
COPY --from=builder /workspace/manager .
COPY --from=builder /workspace/config.ini .

ENTRYPOINT ["/manager"]

