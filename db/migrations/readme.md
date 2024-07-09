# How to run migration

## Run Migration

```go
go run ./db/migrations/migration.go ./db/migrations "host=localhost port=5432 user=postgres dbname=shopeefun-product password=12345 sslmode=disable" up
```

## Down Migration

```go
go run ./db/migrations/migration.go ./db/migrations "host=localhost port=5432 user=postgres dbname=shopeefun-product password=12345 sslmode=disable" down
```

## Create new SQL

```go
go run ./db/migrations/migration.go ./db/migrations "host=localhost port=5432 user=postgres dbname=shopeefun-product sslmode=disable" create add_user_table sql
```
