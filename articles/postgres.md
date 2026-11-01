---
name: PostgreSQL
keywords: [postgres, postgresql, sql, database, db, relational, records, tables, migrations, votes]
status: draft
---
- One database per service, provisioned by the platform. Connection string in the
  environment variable `DATABASE_URL`.
- Use `pgx` (github.com/jackc/pgx/v5). No ORM.
- Schema changes as numbered SQL migrations in `migrations/`, applied at startup.
- No other service reads your tables. Share data through an API.
