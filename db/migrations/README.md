# Database Migrations

This project uses [goose](https://github.com/pressly/goose) for database migrations.

## Automatic Migrations

**Migrations run automatically** when the application starts via `db.ConnectSupabase()`. The migrations are embedded in the binary using `go:embed`, so no external migration files are needed in production.

Simply start your application and migrations will run:

```bash
# Using 1Password CLI
op run --env-file=.env.op -- go run main.go

# Or with environment variables set
go run main.go
```

## Migration Files

Migrations are SQL files in this directory following the goose naming convention:

```
00001_initial_schema.sql
00002_migrate_data.sql
```

Each migration must have:
- `-- +goose Up` directive at the beginning (migration to apply)
- `-- +goose Down` directive at the end (rollback instructions)

## Manual Migration Operations

If you need to manually check migration status or perform rollbacks, you can use the goose CLI directly:

### Install goose CLI (optional)

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Common Commands

```bash
# Check migration status
goose -dir db/migrations postgres "$POSTGRES_URL" status

# Get current version
goose -dir db/migrations postgres "$POSTGRES_URL" version

# Rollback last migration
goose -dir db/migrations postgres "$POSTGRES_URL" down

# Rollback to specific version
goose -dir db/migrations postgres "$POSTGRES_URL" down-to 1

# Create new migration
goose -dir db/migrations create add_new_table sql
```

### With 1Password CLI

```bash
# Check migration status
op run --env-file=.env.op -- goose -dir db/migrations postgres "$POSTGRES_URL" status
```

## Creating New Migrations

1. Create a new migration file:
   ```bash
   go run cmd/migrate/main.go -command=create add_feature_name
   ```

2. Edit the generated file in `db/migrations/`:
   ```sql
   -- +goose Up
   -- Your migration SQL here
   CREATE TABLE new_table (...);

   -- +goose Down
   -- Rollback SQL here
   DROP TABLE new_table;
   ```

3. Test the migration:
   ```bash
   # Apply it
   go run cmd/migrate/main.go -command=up

   # Verify it worked
   go run cmd/migrate/main.go -command=status

   # Test rollback
   go run cmd/migrate/main.go -command=down

   # Reapply
   go run cmd/migrate/main.go -command=up
   ```

## Migration Best Practices

1. **Always write Down migrations** - Make migrations reversible
2. **Test rollbacks** - Verify `down` works before merging
3. **Keep migrations small** - One logical change per migration
4. **Don't modify existing migrations** - Create new ones instead
5. **Use transactions** - Goose wraps migrations in transactions automatically
6. **Idempotent migrations** - Use `IF NOT EXISTS`, `IF EXISTS`, etc.

## Migration Tracking

Goose creates a `goose_db_version` table to track which migrations have run:

```sql
SELECT * FROM goose_db_version ORDER BY id;
```

## Embedded Migrations

Migrations are embedded at compile time:

```go
//go:embed migrations/*.sql
var embedMigrations embed.FS
```

This means:
- ✅ No need to deploy migration files separately
- ✅ Migrations are always in sync with the binary
- ✅ Automatic execution on startup
- ⚠️ Must recompile to pick up new migrations

## Troubleshooting

**Migrations won't run:**
- Check database credentials are set correctly
- Verify `goose_db_version` table exists
- Check migration file naming (must be `00001_name.sql`)

**Migration failed mid-way:**
- Goose uses transactions, so partial failures should rollback
- Check `goose_db_version` table for current state
- May need to manually fix and re-run

**Need to skip a migration:**
```bash
# Mark migration as run without executing it
psql $POSTGRES_URL -c "INSERT INTO goose_db_version (version_id, is_applied) VALUES (2, true);"
```

## References

- [Goose Documentation](https://github.com/pressly/goose)
- [Migration Examples](https://github.com/pressly/goose/tree/master/examples)
