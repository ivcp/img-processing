
FROM golang:1.21.3 AS build-stage

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY ./cmd/api ./api

COPY ./internal ./internal

RUN cd api && go build -o main

FROM gcr.io/distroless/base-debian12 

WORKDIR /

COPY ./migrations ./migrations

COPY --from=build-stage /app/api/main /main

EXPOSE ${SERVER_PORT}

CMD ["/main"] 





