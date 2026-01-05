FROM golang:alpine AS build

WORKDIR /app

RUN apk add --no-cache make gcc musl-dev

COPY go.mod go.sum Makefile .

RUN make prepare

COPY . .

RUN make build

FROM alpine

WORKDIR /app

COPY --from=build /app/backend /backend

ENTRYPOINT ["/backend"]
