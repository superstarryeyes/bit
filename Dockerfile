FROM golang:1.25.0-alpine AS go_builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bit cmd/bit/main.go

CMD ["tail", "-f", "/dev/null"]

FROM alpine:3.22 AS bit

WORKDIR /app

RUN adduser -H -D -u 1000 -s /sbin/nologin app

USER app

COPY --from=go_builder /app/bit /app/bit

ENTRYPOINT [ "/app/bit" ]