---
name: Building a web app or API
keywords: [web, webapp, web app, website, webpage, page, ui, frontend, api, http, server, service, html]
status: draft
---
A web app is one Go service that serves both its API and its static files.

```
myapp/
  go.mod
  main.go          # net/http server, routes, go:embed for static/
  static/
    index.html
    style.css
    app.js
  Dockerfile       # see: goldenpath lookup deployment
```

Rules:
- Embed the frontend with `//go:embed static` so the service is a single binary.
- Listen on the port in the `PORT` environment variable, default `8080`.
- Serve `GET /healthz` returning `200 ok` (the platform's probes use it).
- JSON APIs live under `/api/`, e.g. `GET /api/votes`, `POST /api/votes`.
- Keep state in memory only for prototypes and demos. Real data: `goldenpath lookup storage`.
- Content from outside the company (images, text, APIs): `goldenpath lookup external-services`.
- Log with `log/slog`: `goldenpath lookup observability`.
