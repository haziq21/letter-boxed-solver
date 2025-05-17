ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -v -o /run-app ./cmd/serve/main.go

FROM tiesv/tired-proxy:latest AS tired_proxy

FROM debian:bookworm

COPY --from=builder /run-app /usr/local/bin/
COPY --from=tired_proxy /tired-proxy /usr/local/bin/tired-proxy
RUN apt-get update && apt-get install -y ca-certificates
CMD /usr/local/bin/run-app & \
    /usr/local/bin/tired-proxy --origin=http://localhost:3000 --idle-time 10
