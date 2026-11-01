---
name: Approved external services
keywords: [approved, allowlist, allowed, external, third party, services, apis, cats, cat pictures, cat images, jokes, programming jokes, bad jokes]
status: draft
---
Approved for use from our backends (see `goldenpath lookup external-services` for how):

| Need | Service | How |
|---|---|---|
| Cat pictures | Cat as a service, `https://cataas.com` | `GET https://cataas.com/cat?width=600` returns an image. Proxy it through your backend, e.g. as `GET /api/cat`. |
| Programming jokes | JokeAPI, `https://v2.jokeapi.dev` | `GET https://v2.jokeapi.dev/joke/Programming?safe-mode&type=single` returns JSON with a `joke` field. Always use `safe-mode`. |

Nothing else is approved. Not on the list: ask the platform team.
