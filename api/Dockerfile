FROM golang:alpine AS builder

RUN apk add --no-cache gcc musl-dev sqlite-dev git

WORKDIR /app

COPY go.mod go.sum ./

ENV GOTOOLCHAIN=auto
RUN go mod download

COPY . .

RUN go mod tidy

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite

WORKDIR /root/

COPY --from=builder /app/main .

RUN mkdir -p /root/database

EXPOSE 5000

CMD ["./main"]
