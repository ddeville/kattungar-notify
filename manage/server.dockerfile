FROM golang:1.23.4-bookworm@sha256:2e838582004fab0931693a3a84743ceccfbfeeafa8187e87291a1afea457ff7a AS builder

WORKDIR /usr/src/kattungar-notify

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o exe ./cmd/server

##############

FROM debian:bookworm-20241223-slim@sha256:d365f4920711a9074c4bcd178e8f457ee59250426441ab2a5f8106ed8fe948eb AS service

RUN apt-get update && \
    apt-get install -y ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

COPY --from=builder /usr/src/kattungar-notify/exe /usr/local/bin/kattungar-notify

EXPOSE 5000

CMD ["/usr/local/bin/kattungar-notify"]
