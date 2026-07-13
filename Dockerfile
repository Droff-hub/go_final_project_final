FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o scheduler ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/scheduler .
COPY --from=builder /app/web ./web
EXPOSE 7540
ENV TODO_PORT=7540
CMD ["./scheduler"]