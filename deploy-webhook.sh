#!/bin/bash

echo '---> generating ssl key/cert for webhook if needed: ...'
if ! kubectl get secret admission-webhook-example-certs -o name 2> /dev/null; then
  /webhook-create-signed-cert.sh \
    --service workshopnamespacevalidator \
    --namespace default \
    --secret admission-webhook-example-certs
fi

echo '---> create webhook deployment and svc'
kubectl create -f https://raw.githubusercontent.com/lalyos/k8s-ns-admission/master/deploy.yaml

export CABUNDLE=$(base64 /var/run/secrets/kubernetes.io/serviceaccount/ca.crt | tr -d '\n')

echo '---> create ValidatingWebhookConfiguration'
curl https://raw.githubusercontent.com/lalyos/k8s-ns-admission/master/validating-webh-conf.yaml \
 | envsubst \
 | kubectl apply -f -

echo '---> wait for webhook availability: ...'
while ! curl -k https://workshopnamespacevalidator -m 1 &>/dev/null ; do
  echo -n .
  sleep 1
done

echo '---> testing if the validating webhook works'
kubectl create  ns delme
curl -k https://workshopnamespacevalidator/last 
kubectl delete ns delme