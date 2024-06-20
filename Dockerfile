FROM golang:1.22.2-alpine AS builder

RUN apk update && apk add --no-cache git
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/app


FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
COPY .env .
CMD ./main