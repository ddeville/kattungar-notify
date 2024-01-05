FROM golang:1.21

WORKDIR /usr/src/kattungar-notify

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/kattungar-notify .

EXPOSE 5000

CMD ["kattungar-notify"]
