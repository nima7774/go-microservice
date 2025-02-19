FROM golang:1.23-alpine AS build 

RUN apk --no-cache add gcc g++ ca-certificates

WORKDIR /go/src/github.com/nima7774/go-microservice

COPY go.mod go.sum ./

COPY vandor vendor

COPY account account

COPY catalog catalog

COPY order order

COPY graphql graphql

RUN go build -mod=vendor -o /go/bin/app ./graphql

FROM alpine

WORKDIR /usr/bin

COPY --from=build /go/bin .

EXPOSE 8080

CMD ["./app"]



