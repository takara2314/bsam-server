FROM golang:1.19.3-alpine

WORKDIR /app

COPY . .

ENV GIN_MODE=release
ENV MANAGE_SITE=https://manage.bsam.app

RUN go build -o main .

EXPOSE 8080

CMD ["/app/main"]
