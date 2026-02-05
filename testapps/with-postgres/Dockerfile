

FROM golang:1.25.6-alpine AS builder
WORKDIR /app
 COPY go.mod go.sum ./ 
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/with-postgres .

FROM alpine:latest
COPY --from=builder /app/with-postgres /with-postgres
RUN chmod +x /with-postgres
ENTRYPOINT ["/with-postgres"]
