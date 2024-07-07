FROM golang:1.22.1-alpine


WORKDIR /app

COPY /../.env .

COPY ./scripts .

RUN go build -o init_rabbitmq init_rabbitmq.go

ENTRYPOINT ["./init_rabbitmq"]