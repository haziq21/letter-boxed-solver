ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go build -v -o /run-app ./cmd/serve/main.go


FROM debian:bookworm

COPY --from=builder /run-app /usr/local/bin/
RUN apt-get update && apt-get install -y ca-certificates
CMD ["run-app"]
