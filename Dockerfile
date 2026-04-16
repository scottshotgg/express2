FROM golang:1.25-alpine AS builder
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /express2 ./...

FROM alpine:3.19
COPY --from=builder /express2 /usr/local/bin/express2
ENTRYPOINT ["/usr/local/bin/express2"]
