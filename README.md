# Table of Contents

- [Local Environment Setup](#local-environment-setup)
  - [Postgres](#postgres)
  - [Migration](#migration)
  - [Hot Reload](#hot-reload)
- [Useful Commands](#useful-commands)
  - [pgcli](#pgcli)
  - [iredis](#iredis)

## Local Environment Setup

### Postgres

```bash
# Setup local docker container
docker pull postgres:15.4
docker run -d -e POSTGRES_USER=everytrack -e POSTGRES_PASSWORD=everytrack -p 5432:5432 -v /var/lib/postgresql/data/everytrack:/var/lib/postgresql/data --name everytrack-pg postgres:15.4

# Access database in container
docker exec -it everytrack-pg bash
psql -U everytrack -d postgres
```

### Migration

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Export path to global golang executable bin directory
vim .zshrc
export PATH=$PATH:$HOME/go/bin

# Create new migration file
migrate create -ext sql -dir ./migrations -seq <migration_file_name>

# Run migration
migrate -path ./migrations -database "postgresql://everytrack:everytrack@127.0.0.1:5432/everytrack?sslmode=disable" up

# Rollback migration
migrate -path ./migrations -database "postgresql://everytrack:everytrack@127.0.0.1:5432/everytrack?sslmode=disable" down -all
```

### Hot Reload

```bash
go install github.com/cosmtrek/air@latest

# Run the server
air
```

## Useful Commands

### pgcli

We will be using a powerful cli tools to manage our postgres database - `pgcli`

```bash
brew install pgcli
vim ~/.config/pgcli/config

# A config file is automatically created at ~/.config/pgcli/config at first launch
# See the file itself for a description of all available options
# Add alias dsn config under the section like below
[alias_dsn]
# example_dsn = postgresql://[user[:password]@][netloc][:port][/dbname]
pgcli -D <name>
```

### iredis

We will be using a powerful cli tools to manage our redis - `iredis`

```bash
brew install iredis
vim ~/.iredisrc

# Add alias dsn config inside ~/.iredisrc
[alias_dsn]
<name>=redis://<username>:<password>@<host>:<port>
iredis -d <name>
```
