FROM golang:1.11-alpine
RUN apk add -U git musl-dev gcc
ADD . /app
WORKDIR /app
RUN go build -v -o webhook_linux .

FROM alpine:latest
COPY --from=0 /app/webhook_linux /webhook
ENTRYPOINT ["/webhook"]