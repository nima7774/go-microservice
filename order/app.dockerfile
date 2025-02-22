FROM golang:1.22-alpine3.18 AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/nima7774/go-microservice
COPY go.mod go.sum ./
COPY vendor vendor
COPY account account
COPY catalog catalog
COPY order order
RUN go build -mod=vendor -o /go/bin/app ./order/cmd/order 

FROM alpine:3.18
WORKDIR /usr/bin
COPY --from=build /go/bin .
EXPOSE 8080
CMD ["app"]

