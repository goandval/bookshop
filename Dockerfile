# syntax=docker/dockerfile:1
FROM golang:1.24-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bookshop ./cmd/bookshop

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /bookshop /bookshop
COPY configs ./configs
COPY migrations ./migrations
ENV GIN_MODE=release
CMD ["/bookshop"] 