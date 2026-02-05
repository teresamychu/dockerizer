

FROM golang:1.25.6-alpine AS builder
WORKDIR /app
 COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /app/simple .

FROM alpine:latest
COPY --from=builder /app/simple /simple
RUN chmod +x /simple
ENTRYPOINT ["/simple"]
