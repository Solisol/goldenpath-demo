---
name: Supported languages and frameworks
keywords: [language, languages, go, golang, backend, frontend, html, css, javascript, typescript, react, framework]
status: draft
---
- **Backend: Go**, standard library first. `net/http` is enough for most services; add a
  dependency only when the standard library clearly falls short.
- **Frontend, small apps: plain HTML, CSS and JavaScript**, served by the Go service. No
  build step, no framework, no npm.
- **Frontend, larger apps: React with TypeScript.** Talk to the platform team first.

Anything else (Python, Rust, Node backends, other frontend frameworks) needs an ADR before
the first line of code. A new language spreads into CI, deploys, on-call and hiring.

How to structure a web app: `goldenpath lookup web-app`.
