
FROM golang:1.21.3

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd/api ./api

COPY ./internal ./internal

RUN cd api && go build -o main

EXPOSE ${SERVER_PORT}

CMD cd api && ./main 



