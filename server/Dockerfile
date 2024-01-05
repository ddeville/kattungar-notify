FROM golang:1.21 AS builder

WORKDIR /usr/src/kattungar-notify

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o exe .

##############

FROM ubuntu:22.04

COPY --from=builder /usr/src/kattungar-notify/exe /usr/local/bin/kattungar-notify

EXPOSE 5000

CMD ["/usr/local/bin/kattungar-notify"]
