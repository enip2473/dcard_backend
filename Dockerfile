FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/src/cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -tags=viper_bind_struct -o /server main.go

FROM gcr.io/distroless/base-debian10
COPY --from=builder /server /
CMD ["/server"]
