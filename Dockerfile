FROM golang:alpine AS builder
WORKDIR /server
COPY . .
RUN go mod tidy
RUN go build -o main cmd/main.go

FROM alpine
WORKDIR /server
COPY --from=builder /server/main /server/main
CMD ["./main"]
