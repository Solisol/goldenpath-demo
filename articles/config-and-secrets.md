---
name: Configuration and secrets
keywords: [config, configuration, environment, env, env vars, secrets, secret, password, api key, token, credentials]
status: draft
---
- All configuration through environment variables, with safe defaults for local runs.
- Secrets come from the platform's secret store as environment variables at runtime.
  Never in code, config files, Dockerfiles or git, not even "just for the demo".
- Document every variable the service reads in its README.
