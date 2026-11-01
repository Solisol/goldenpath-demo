---
name: Running on Kubernetes
keywords: [kubernetes, k8s, k3s, cluster, namespace, deploy, deployment, run, hosting, production, container, docker, dockerfile, image]
status: draft
---
- **Backend services run on Kubernetes.** For now that's the internal demo cluster.
- **Each product gets its own namespace**, named after the product (e.g. `catvote`).
  Nothing runs in `default`.
- One container image per service: multi-stage Dockerfile, `golang` to build,
  `gcr.io/distroless/static` to run, numeric non-root user.
- Namespaces enforce the *restricted* Pod Security Standard: non-root, no privilege
  escalation, all capabilities dropped, read-only root filesystem.
- Every service needs `GET /healthz` and logs to stdout.
- Pin image tags (never `:latest`) and set memory limits.
