# goldenpath-demo

A small golden-path service for humans and coding agents, from the talk
*Practical guide to agent native platforms* (Øredev 2026).

A platform's golden path ("how we build things here") is usually a wiki that people
skim once. Coding agents don't skim; they need the right piece of it at the moment
they're about to make a choice. This demo exposes the golden path as an API with a
thin command-line client, so a developer or an agent can ask "how do we do storage
here?" just in time.

```
$ goldenpath search storage
storage  Storing data  (matched: storage)
buckets  Object storage (buckets)  (matched: object storage, name)

$ goldenpath lookup buckets
# Object storage (buckets)
- One bucket per service, created by the platform. The name comes from ...
```

## Why a thin client and a service

Whatever you install on developer laptops is hard to change: you don't know which
version people run, and you can't update it for them. So the `goldenpath` CLI does
nothing but forward the question over HTTP and print the answer. Everything that
changes (articles, search, guidelines) lives in the service, which the platform team
controls and can update at any time. The CLI sends its version with every request,
so the platform team can see which clients are out there.

The service logs every search and lookup and shows them live at `/`. For the platform
team, that's product telemetry: what do developers and agents ask for, and what don't
they find?

## What's in here

```
articles/                the golden path: one Markdown file per article
articles.go              embeds the articles into the service
cmd/goldenpath/          the CLI
cmd/goldenpath-api/      the HTTP service
deploy/goldenpath.yaml   Kubernetes manifests (namespace idp, Deployment, Service, Ingress)
Dockerfile, Makefile
```

The articles are examples for a made-up company. Replace them with your own.

## Articles

Each article is a Markdown file. The file name is the article id.

```markdown
---
name: Object storage (buckets)
keywords: [bucket, object storage, s3, files, images]
---
Markdown text. Refer to other articles with the command that fetches them,
e.g. `goldenpath lookup external-services`, so an agent can follow the trail.
```

The articles are built into the service. To serve them from a directory instead
(a wiki export, a git checkout), set `GOLDENPATH_CONTENT_DIR`.

## API

| Endpoint | Returns |
|---|---|
| `GET /search?q=<words>` | Articles whose keywords or name partially match, best match first: `{"query": ..., "results": [{"id", "name", "keywords", "matched"}]}` |
| `GET /lookup/<id>` | One article: `{"id", "name", "keywords", "content"}`, or 404 |
| `GET /` | Live activity page |
| `GET /activity` | The same activity as server-sent events |
| `GET /healthz` | `ok` |

## Try it

You need Go 1.27 or later.

Install the CLI (the service URL defaults to `http://localhost:8080`; change it with
`GOLDENPATH_URL` or `--url`):

```sh
go install github.com/Solisol/goldenpath-demo/cmd/goldenpath@latest
```

Run the service, pick one:

```sh
# straight from source
go run ./cmd/goldenpath-api

# as a container
docker build -t goldenpath-api:0.1.0 .
docker run --rm -p 8080:8080 goldenpath-api:0.1.0

# in a local k3s cluster (k3d), reachable on localhost:8080
k3d cluster create goldenpath -p "8080:80@loadbalancer" --wait
docker build -t goldenpath-api:0.1.0 .
k3d image import goldenpath-api:0.1.0 -c goldenpath
kubectl apply -f deploy/goldenpath.yaml
```

The manifests put the service in its own namespace, `idp`, which enforces the
Kubernetes *restricted* Pod Security Standard: non-root, no privilege escalation,
all capabilities dropped, read-only root filesystem, no service-account token.

Then:

```sh
goldenpath search cat pictures
goldenpath lookup approved-services
open http://localhost:8080          # watch the lookups arrive
```

The `Makefile` has the same steps: `make test`, `make cli`, `make cluster`, `make deploy`.

## Using it with a coding agent

Put something like this in the repo's `AGENTS.md` (or `CLAUDE.md`):

```markdown
Our golden path lives in the platform. Before choosing a technology, library or
pattern, run `goldenpath search <words>`, then `goldenpath lookup <id>`, and follow
the references until you know how we do it here.
```

A small always-present pointer like this, with the details fetched on demand, keeps
the agent's context small without leaving it to guess when to look things up.

## License

MIT
