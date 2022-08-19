FROM golang:1.18.3 as cache
WORKDIR /app
COPY go.* ./
RUN go mod download

# Test stage
FROM cache as image-test
COPY . ./
ENV GIN_MODE=release

# Prod/Build stage
FROM cache as builder
COPY . ./
RUN CGO_ENABLED=0 go build -o dionysos-server

FROM alpine:3.16.0 as image-prod
COPY --from=builder /app/dionysos-server /
EXPOSE 8080
CMD ["/dionysos-server"]
