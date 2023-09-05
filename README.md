# everytrack-backend

## Local Environment Setup

### Migration

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Export path to global golang executable bin directory
vim .zshrc
export PATH=$PATH:$HOME/go/bin

# Create new migration file
migrate create -ext sql -dir ./migrations -seq <migration_file_name>
```
