FROM golang:1.18-alpine as cache
WORKDIR /app
COPY go.* .
RUN go mod download

FROM cache as builder
COPY . .
RUN CGO_ENABLED=0 go build -o dionysos-server

# Prod stage
FROM alpine:3.16.0 as image-prod
COPY --from=builder /app/dionysos-server /
EXPOSE 8080
CMD ["/dionysos-server"]
