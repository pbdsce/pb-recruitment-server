# How to run?
- Install [air](https://github.com/air-verse/air?tab=readme-ov-file#installation)
- Run `air` command in the terminal
- Server is running on `http://localhost:8080`

## DB Migrations

- Install [migrate](https://github.com/golang-migrate/migrate)
- Add `DB_ADDR` to `.env`

### Migration Commands

```bash
# Create a new migration
make migration create_users_table

# Run migrations up
make migrate-up

# Rollback migrations (specify number of steps)
make migrate-down 1

# Force migration version (use with caution)
make migrate-force 1
```