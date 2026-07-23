FROM golang:latest as builder


WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 go build -o grpc-jwt ./cmd/grpc-jwt/main.go

FROM scratch

WORKDIR /app

COPY --from=builder /app/grpc-jwt ./

COPY --from=builder /app/.env ./

COPY --from=builder /app/certs ./certs

EXPOSE 50052

ENTRYPOINT [ "./grpc-jwt" ]