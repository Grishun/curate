FROM golang:1.25-alpine3.22 AS builder

WORKDIR /app

COPY . ./

RUN go mod download


RUN go build -o /service cmd/service/main.go

FROM alpine:3.22

COPY --from=builder /service /service

CMD ["/service"]

EXPOSE 8080
