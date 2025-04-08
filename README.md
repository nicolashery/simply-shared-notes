# Simply Shared Notes

**Simply Shared Notes** lets you create and share notes with others — no accounts, no passwords, no friction.
Create a **Space**, add one or more notes in plain text Markdown, and share the **Access Link** to give others instant access.

Each Space has one or more **Members**, allowing people to identify themselves and have their contributions tracked — all without needing to register. It’s built on trust, designed for small private groups.

## Features

- ⚡️ **No signup or login required**
- 🔑 **Multiple Access Link types**
  - **Admin** – full control
  - **Edit** – can add/edit notes
  - **View** – read-only access
- 📨 **Recover lost Access Links via email**
- 👥 **Named Members for each Space**
  - Visitors select their identity from the Member list
  - Their choice is remembered on the device
  - Notes show who created or last updated them
  - View-only links skip Member selection
- 📚 **Activity history for each Space**
  - Tracks who did what and when
  - Create, edit, and delete actions are recorded per Member
- 📝 **Plain text notes in Markdown**
- 📱 **Mobile-friendly**
  - Minimal JavaScript, works on all devices

## Tech stack

- [Go](https://go.dev/): Programming language
- [Chi](https://go-chi.io): Router
- [Templ](https://templ.guide): HTML templates
- [SQLite](https://www.sqlite.org/): Database
- [sqlc](https://github.com/sqlc-dev/sqlc): Generate code from SQL
- [Dbmate](https://github.com/amacneil/dbmate): Database migrations
- [Tailwind](https://tailwindcss.com/) and [daisyUI](https://daisyui.com/): CSS
- [htmx](https://htmx.org/): Client-side interactions
- [Vite](https://vite.dev/): Assets tooling
- [Fly.io](https://fly.io/): Deployment

## Development

Install Go, for example using [mise](https://mise.jdx.dev/lang/go.html):

```bash
mise use -g go@latest
```

Install tools:

```bash
brew install go-task/tap/go-task
go install github.com/a-h/templ/cmd/templ@latest
```

Run the application:

```bash
task run
```

For a full list of tasks run `task -a`.
