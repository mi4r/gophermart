FROM golang:1.22.5-alpine AS builder

WORKDIR /app

COPY . .

ARG SERVICE_NAME
RUN go build -o ${SERVICE_NAME} cmd/${SERVICE_NAME}/main.go

FROM alpine:3.19

WORKDIR /app

ARG SERVICE_NAME
COPY --from=builder /app/internal/storage/default/migrations /app/internal/storage/default/migrations
COPY --from=builder /app/${SERVICE_NAME} .
CMD ["./${SERVICE_NAME}"]