
[![Docker Automated build](https://img.shields.io/docker/automated/lalyos/k8s-ns-admission.svg)](https://hub.docker.com/r/lalyos/k8s-ns-admission/)

# Namespace Validaing Admission Webhook

This is a Validating Admission Webhook for k8s. It is validating namespaces create/delete operations, checking name prefixes. The intent is to let the `user1` ServiceAccount be able to create/delete any namespace with names starting with `user1-`

## Usage

To deploy the validating admission webhook, just run this one-liner:
```
kubectl apply -f https://raw.githubusercontent.com/lalyos/k8s-ns-admission/master/deploy-webhook-job.yaml 
```

Most admission webhook examples requires you to do a lot of manual steps regarding certificate creation. That is all taken care by the job:

- creates a service account to give certificate relatades access to the job
- creates a certificate req for the webhook's https endpoint
- let api-server sign the csr
- save the cert/key into a secret
- create a svc/deployment for the webhook (using the generated cert as mounted volume)
- creates the ValidatingWebhookConfiguration
- tests the webhook by creating and deleteing a ns "delme"

## tl;dr

For k8s workshops you might want to provide isolated workspace for each participant.
Let's say you create a **separate namespace** for each of them. To restrict them to only that specific namespace, we can use RBAC. You can create a ServiceAccount with an appropriate Role to give full access to namespaced resources (pod,deployment,svc,job ...)

But how can they learn to work with namespaces? We might want to give access to `user1` to create/delete any namespace with names starting with `user1-`. Unfortunately RBAC Roles doesnt support pattern based roles. For details see: [k8s#56582](https://github.com/kubernetes/kubernetes/issues/56582)

Then you might think, lets just give them 1 single predefined ns to create, like `user1-play`. Therefoe lets create a ClusterRole which gives access to namespaces resources restricted by `resourceNames: ["user1-play"]`. Unfortunately its not possible, as there is a sidenote in [RBAC docs](https://kubernetes.io/docs/reference/access-authn-authz/rbac/#referring-to-resources)
> Notably, if resourceNames are set, then the verb must not be list, watch, create, or delete

## References

code is proudly stolen from: 
- [official k8s webhook test](https://github.com/kubernetes/kubernetes/tree/v1.10.0-beta.1/test/images/webhook)
- [istio cert-gen script](https://github.com/istio/istio/raw/release-0.7/install/kubernetes/webhook-create-signed-cert.sh)

relevant docs:
- [Dynamic Admission Control](https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/)
