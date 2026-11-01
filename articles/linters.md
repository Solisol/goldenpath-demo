---
name: Linters and fast feedback
keywords: [lint, linter, linters, linting, check, checks, rules, feedback, ci, quality, gofmt, vet, failing, failed]
status: draft
---
We use linters to enforce good practices in fast feedback loops: the earlier a mistake
is caught, the cheaper it is to fix. Every repo has the platform linter in `lint/`.

- Run it: `./lint/lint.sh`. For coding agents it also runs automatically after every
  edit, and the result comes back right away.
- **It is forbidden to turn off linters, or to change them to make a check pass.**
  They are there for a reason. If a rule seems wrong, ask the platform team.
- When a check fails, its message names the rule and the golden-path article that
  explains it. Read the article, then fix the code.

What it checks today: Go formatting and `go vet`, no Node.js, Go only for backends,
nothing loaded from third parties in the browser, no hardcoded secrets, and that the
linter itself hasn't been changed.
