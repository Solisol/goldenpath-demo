---
name: Frontend tooling and libraries
keywords: [node, nodejs, node.js, npm, yarn, pnpm, package.json, build tool, bundler, vite, frontend dependencies, javascript library, javascript libraries]
status: draft
---
We don't use Node.js anywhere: no npm, no `package.json`, no bundlers or build steps.
Every Node toolchain is one more thing to patch, audit and keep running, for little
gain in small apps.

Instead:
- Write it yourself in plain JavaScript. Modern browsers cover most needs (`fetch`,
  `<canvas>`, CSS animations, `<dialog>`, ES modules).
- Need a small library anyway? Vendor one file into `static/vendor/` with its licence
  header intact, and serve it from our own service.
- Never load scripts, styles or fonts from a CDN: `goldenpath lookup external-services`.

The linter fails the build if Node files appear: `goldenpath lookup linters`.
