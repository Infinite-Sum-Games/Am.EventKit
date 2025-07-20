# 🛠️ Anokha Backend

Backend server for Anokha, built with Go, PostgreSQL, Redis, and modern tooling.

---

## 📦 Key Libraries

- **Gin** – Fast HTTP web framework  
- **sqlc** – Type-safe Postgres query generation  
- **PostgreSQL** – Relational DB  
- **Redis** – Caching & session management  
- **zerolog** – High-performance structured logging  
- **Viper** – Configuration via YAML/env  
- **ozzo-validator** – Declarative validation  
- **Lefthook** – Git hooks for code quality  
- **Docker** – Local development & deployment

---

## 🚀 Getting Started

### 1. Install Dependencies

You can either install them manually or use Docker.

#### ✅ Go (version 1.24)

- **macOS:** `brew install go@1.24`  
- **Ubuntu/Debian:** `sudo apt update && sudo apt install golang`  
- **Windows:** Use [official installer](https://go.dev/dl)

### 🧩 Option A: Manual Install (OS-Specific)

#### ✅ PostgreSQL

- **macOS:**  
  ```bash
  brew install postgresql
  brew services start postgresql
  ```

- **Ubuntu/Debian:**  
  ```bash
  sudo apt update && sudo apt install postgresql postgresql-contrib
  sudo systemctl start postgresql
  ```

- **Windows:** Install from https://www.postgresql.org/download/windows/

#### ✅ Goose (for DB migrations)

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Make sure `$GOPATH/bin` is in your `$PATH`.

#### ✅ Lefthook (optional, for git hooks)

```bash
go install github.com/evilmartians/lefthook@latest
lefthook install
```

### 🐳 Option B: Use Docker (Recommended)

You can run Postgres and Redis with Docker:

```bash
docker compose up -d
```

This will start:

- PostgreSQL on `localhost:5432`
- Valkey on `localhost:6379`
- RedisInsight on `localhost:5540`

---

### 2. Configure Environment

Create an `env.toml` file with the appropriate values (refer to `env.sample.toml`).

---

### 🛠️ Install `make` (if not installed)

#### ✅ macOS
```bash
brew install make
```

#### ✅ Ubuntu/Debian
```bash
sudo apt update && sudo apt install build-essential
```

#### Windows
- Recommended: Install **Make** via Chocolatey:
  ```bash
  choco install make
  ```
- Or use [GNUWin Make](http://gnuwin32.sourceforge.net/packages/make.htm) and add it to your system `PATH`.

---

### 3. Initialize Project Dependencies

Run the following command to install all required CLI tools and set up the project:

```bash
make setup
```
---

### 4. Run the Server

Use the included `Makefile`:

```bash
make run     # Start the Go server

make doc     # (optional) Start services with Docker
make up      # Run DB migrations
```

Other useful commands:

```bash
make seed     # Seed the database
make test     # Run all tests
make down     # Roll back migrations
```

---

## 🧪 Testing

```bash
make test
```

---

## 🤝 Contributing

1. Create a feature branch
2. Push your changes and open a PR

### 🌿 Branch Naming Convention

Use the following prefixes to name your branches based on the type of work:

| Prefix     | Use for                       | Example                    |
|------------|-------------------------------|----------------------------|
| `feat/`    | New features                  | `feat/user-login`          |
| `fix/`     | Bug fixes                     | `fix/event-date-validation`|
| `docs/`    | Documentation changes         | `docs/backend-api`         |
| `hotfix/`  | Urgent production patches     | `hotfix/fix-email-crash`   |
| `chore/`   | Misc tasks (e.g., deps, CI/CD)| `chore/update-modules`     |

> ✅ Always use descriptive names after the prefix for clarity.
