# SQL

## Schema

The database schema is automatically dumped by [Dbmate](https://github.com/amacneil/dbmate) and can be found in the file [sql/schema.sql](../sql/schema.sql).

Common columns:

- `id`
  - Primary key, auto-incremented integer
  - For internal use only (ex: foreign keys, table joins)
  - We don't use SQLite's [AUTOINCREMENT](https://sqlite.org/autoinc.html), as recommended for efficiency reasons
- `created_at`/`updated_at`
  - UTC datetime
  - Set by application code (not SQL schema or triggers)
- `created_by`/`updated_by`
  - Foreign key to `members` table, nullable
  - Initially NULL for a new entry in `spaces` and the first entry in `members` for that space
  - After creating both the space and first member in a transaction, the application updates these fields to reference the first member
  - For subsequent operations, these fields reference the member who performed the action

## Public IDs and Access Tokens

Internal IDs:

- Integer
- All tables
- Column: `id`
- Primary key
- For internal use only, will never appear in a URL

Access Tokens:

- Random string: alphabet `A-Za-z0-9`, length `20`
- `spaces` table
- Columns: `admin_token`, `edit_token`, `view_token`
- Unique constraint, index
- Used in URLs to access a Space

Public IDs:

- Random string: alphabet `A-Za-z0-9`, length `10`
- Space object tables (ex: `members`, `notes`)
- Column: `public_id`
- Unique constraint on `(space_id, public_id)`, composite index
- Used in URLs, in combination with a Space's Access Token

## Deleting members

- When a Member is deleted, their references in other tables (`created_by` and `updated_by`) are set to `NULL`
- Foreign keys that reference members use `ON DELETE SET NULL`
- Application handles `NULL` values for these references and the UI displays "deleted member" instead of the member name (which is no longer available)

## Conventions

We roughly follow [Ruby on Rail's](https://guides.rubyonrails.org/active_record_basics.html#naming-conventions) schema naming conventions:

- Plural names for tables (ex: `members`)
- `id` for primary key column
- `_id` suffix for foreign keys (ex: `space_id`)
- `idx_<table_name>_<column_name>` for indices (append more columns as necessary separated by `_`)

For source SQL files (migrations, queries):

- Use all caps for SQL keywords (ex: `SELECT * FROM members`)

## Migrations

- [Dbmate](https://github.com/amacneil/dbmate) is the tool used to run SQL migrations
- Run migrations with `task migrate`
- Migration files located under `sql/migrations/`
- No "down" migrations, only "up"
  - If a fix needs to be made in production, and "up" migration is quickly written to patch it
  - For local development, the workflow to go back is to drop the database (`task drop`) and re-run migrations (`task migrate`)
- **Note**: SQLite only supports a few basic alter table operations ("rename table", "rename column", "add column", "drop column")
  - For more complex operations see the [ALTER TABLE documentation](https://sqlite.org/lang_altertable.html)

## Seed data

- Some sample seed data to help with local development is located at `sql/seed.sql`
- We use pure SQL files for seed data for portability and to be able to use a database dump
- Load the seed data on a fresh database after migrations have been run with `task seed`
- Uses the `sqlite3` CLI (present on most systems)
- Not idempotent, does not check if data is already loaded (to be run on fresh database after migrations)
- Workflow to reset database with seed data is: drop database, run migrations, load seed data (shortcut: `task reset`)

## SQLite settings

- To optimize SQLite for concurrency, disk and memory usage, and data integrity, all in the context of a web application, a set of [PRAGMA statements](https://sqlite.org/pragma.html) are used
- The values used are located in `sql/pragmas.sql`
- PRAGMA statements are set by the application when opening the SQLite connection
  - They are not set by a migration (since some if not most PRAGMA statements are valid for the scope of a connection only)

## Code generation

- [sqlc](https://sqlc.dev/) is used to generate Go code from SQL
- SQL queries are located in `sql/queries/`
- sqlc uses the `sql/schema.sql` that is automatically dumped by Dbmate on every migration run
