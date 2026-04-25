FROM golang:1.26.1-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main ./cmd/server/
EXPOSE 8080
CMD ["/app/main"]

#FROM alpine:3.23
#WORKDIR /app
#COPY --from=builder /app/main /app
#EXPOSE 8080
#CMD ["/app/main"]


