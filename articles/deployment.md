---
name: Building and deploying
keywords: [deploy, deployment, ship, release, container, docker, dockerfile, image, kubernetes, k8s, k3s, production]
status: draft
---
- One container image per service, built from a multi-stage Dockerfile:
  `golang` to build, `gcr.io/distroless/static` to run, as a non-root user.
- The platform provides the Kubernetes manifests from a template; you set name, port
  and environment variables. Don't hand-write Ingress or RBAC.
- Every service must pass `/healthz` and log to stdout before it can be deployed.
