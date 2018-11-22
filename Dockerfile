FROM golang:1.11-alpine
RUN apk add -U git musl-dev gcc
ADD . /app
WORKDIR /app
RUN go build -v -o webhook_linux .

FROM alpine:latest
RUN apk add -U curl bash openssl gettext
ADD https://raw.githubusercontent.com/istio/istio/release-0.7/install/kubernetes/webhook-create-signed-cert.sh /
RUN chmod +x webhook-create-signed-cert.sh
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl \
   && chmod +x ./kubectl \
   && mv ./kubectl /usr/local/bin/kubectl
ADD deploy-webhook.sh /

COPY --from=0 /app/webhook_linux /webhook
ENTRYPOINT ["/webhook"]