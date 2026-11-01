---
name: Supported languages
keywords: [language, languages, go, golang, backend, frontend, html, css, javascript, stack, tech stack]
status: draft
---
- **Backend: Go.** Standard library first; `net/http` is enough for most services. Add a
  dependency only when the standard library clearly falls short.
- **Frontend: keep it as slim as possible.** Plain HTML, CSS and JavaScript, served by
  the Go service. See `goldenpath lookup web-app`.

Other backend languages (Python, Java, Rust, ...) need an ADR before the first line of
code. A new language spreads into CI, deploys, on-call and hiring.
