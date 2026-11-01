---
name: Storing data
keywords: [storage, store, data, persist, persistence, database, db, files, state, save]
status: draft
---
Pick by the shape of the data:
- **Files, images, documents, exports** → object storage (buckets): `goldenpath lookup buckets`
- **Records, relations, transactions, anything you query** → PostgreSQL: `goldenpath lookup postgres`
- **Prototypes and demos** → in memory is fine; say so in the README.

We don't run other databases or storage services (no MongoDB, Redis as a database, SQLite
in production, or files on the container's disk). Need one? Write an ADR first.
