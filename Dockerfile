FROM golang:1.20.4-alpine

WORKDIR /app

COPY . .

RUN go build -o main .

EXPOSE 8080

CMD ["/app/main"]
