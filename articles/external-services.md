---
name: Using services outside the company
keywords: [external, third party, third-party, external api, api, integration, internet, fetch, images, pictures, photos, cats, jokes, content, embed, hotlink, licence, license]
status: draft
---
Anything from outside the company (APIs, images, fonts, scripts, text) is a dependency
we have to trust, keep available and be allowed to use.

Rules:
1. **Only use approved services**: `goldenpath lookup approved-services`.
   Not on the list? Ask the platform team; don't add it yourself.
2. **Call them from the backend, never from the browser.** The page talks only to our
   own service. (Privacy: the visitor's browser never contacts a third party.)
3. **Timeout of 2 seconds** on every call, and a **fallback** so the page still works
   when the service is down (a default image or text that ships with the app).
4. **No hotlinking**: don't put third-party URLs in our HTML. Proxy or copy through the backend.
5. Log every failed call with `log/slog` at warn level.
