---
layout: "kubernetes"
page_title: "Kubernetes: kubernetes_secret"
sidebar_current: "docs-kubernetes-resource-secret"
description: |-
  The resource provides mechanisms to inject containers with sensitive information while keeping containers agnostic of Kubernetes.
---

# kubernetes_secret

The resource provides mechanisms to inject containers with sensitive information, such as passwords, while keeping containers agnostic of Kubernetes.
Secrets can be used to store sensitive information either as individual properties or coarse-grained entries like entire files or JSON blobs.

## Example Usage

```
resource "kubernetes_secret" "example" {
  metadata {
  	name = "my-secret"
  }
  data {
  	api_key = "8Rj$FcK9w6"
  	db_password = "63IiH*#7nt"
  }
}
```

## Argument Reference

The following arguments are supported:

* `data` - (Optional) A map of the secret data.
* `metadata` - (Required) Standard secret's metadata. More info: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#metadata
* `type` - (Optional) The secret type, defaults to `Opaque`. Other common types involve `kubernetes.io/tls` which is what gets created when using `kubectl create secret tls` and which can be used for TLS Ingresses. Changing this value forces recreation of the secret. More info: https://kubernetes.io/docs/user-guide/kubectl/kubectl_create_secret_tls/ and https://github.com/kubernetes/ingress/blob/master/controllers/gce/README.md#tls

## Nested Blocks

### `metadata`

#### Arguments

* `annotations` - (Optional) An unstructured key value map stored with the secret that may be used to store arbitrary metadata. More info: http://kubernetes.io/docs/user-guide/annotations
* `generate_name` - (Optional) Prefix, used by the server, to generate a unique name ONLY IF the `name` field has not been provided. This value will also be combined with a unique suffix. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#idempotency
* `labels` - (Optional) Map of string keys and values that can be used to organize and categorize (scope and select) the secret. May match selectors of replication controllers and services. More info: http://kubernetes.io/docs/user-guide/labels
* `name` - (Optional) Name of the secret, must be unique. Cannot be updated. More info: http://kubernetes.io/docs/user-guide/identifiers#names
* `namespace` - (Optional) Namespace defines the space within which name of the secret must be unique.

#### Attributes

* `generation` - A sequence number representing a specific generation of the desired state.
* `resource_version` - An opaque value that represents the internal version of this secret that can be used by clients to determine when secret has changed. Read more: https://github.com/kubernetes/community/blob/master/contributors/devel/api-conventions.md#concurrency-control-and-consistency
* `self_link` - A URL representing this secret.
* `uid` - The unique in time and space value for this secret. More info: http://kubernetes.io/docs/user-guide/identifiers#uids

## Import

Secret can be imported using its name, e.g.

```
$ terraform import kubernetes_secret.example my-secret
```
