FROM golang:1.23-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main

FROM alpine:latest
WORKDIR /
COPY --from=builder /main /main
EXPOSE 23234/tcp
CMD ["/main"]