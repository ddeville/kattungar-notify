FROM golang:1.21-bookworm AS builder

WORKDIR /usr/src/kattungar-notify

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o exe ./cmd/server

##############

FROM debian:bookworm-slim

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/src/kattungar-notify/exe /usr/local/bin/kattungar-notify

EXPOSE 5000

CMD ["/usr/local/bin/kattungar-notify"]
