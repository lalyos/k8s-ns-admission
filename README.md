
# Namespace Validaing Admission Webhook

This is a Validating Admission Webhook for k8s. It is validating namespaces create/delete operations, checking name prefixes. The intent is to let the `user1` ServiceAccount be able to create/delete any namespace with names starting with `user1-`

## tl;dr

For k8s workshops you might want to provide isolated workspace for each participant.
Let's say you create a **separate namespace** for each of them. To restrict them to only that specific namespace, we can use RBAC. You can create a ServiceAccount with an appropriate Role to give full access to namespaced resources (pod,deployment,svc,job ...)

But how can they learn to work with namespaces? We might want to give access to `user1` to create/delete any namespace with names starting with `user1-`. Unfortunately RBAC Roles doesnt support pattern based roles. For details see: (k8s#56582)[https://github.com/kubernetes/kubernetes/issues/56582]

Then you might think, lets just give them 1 single predefined ns to create, like `user1-play`. Therefoe lets create a ClusterRole which gives access to namespaces resources restricted by `resourceNames: ["user1-play"]`. Unfortunately its not possible, as there is a sidenote in (RBAC docs)[https://kubernetes.io/docs/reference/access-authn-authz/rbac/#referring-to-resources]
> Notably, if resourceNames are set, then the verb must not be list, watch, create, or delete

See (Dynamic Admission Control docs:)[https://kubernetes.io/docs/reference/access-authn-authz/extensible-admission-controllers/]

## Origin
Code started from the: (official k8s webhook test)[https://github.com/kubernetes/kubernetes/tree/v1.10.0-beta.1/test/images/webhook]