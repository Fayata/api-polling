FROM golang:1.21.5-alpine

WORKDIR /app

COPY . .

RUN go build -o api-polling

CMD ["./api-polling"]

# CMD ["tail", "-f", "/dev/null"]
