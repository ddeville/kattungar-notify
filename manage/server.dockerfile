FROM golang:1.26.0-trixie@sha256:889885d7cc1275935e3f9920aabadc5fadbe873f633d92a746f1bc401dd40f69 AS builder

WORKDIR /usr/src/kattungar-notify

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o exe ./cmd/server

##############

FROM debian:trixie-20260202-slim@sha256:f6e2cfac5cf956ea044b4bd75e6397b4372ad88fe00908045e9a0d21712ae3ba AS service

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/src/kattungar-notify/exe /usr/local/bin/kattungar-notify

EXPOSE 5000

CMD ["/usr/local/bin/kattungar-notify"]
