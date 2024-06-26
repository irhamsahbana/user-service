### Run Migration
```
go run migration.go ./sql "host=localhost port=5432 user=localhost dbname=db_users password=postgres sslmode=disable" up
```

### Down Migration
```
go run migration.go ./sql "host=localhost port=5432 user=localhost dbname=db_users password=postgres sslmode=disable" down
```

### Create new SQL
```
go run migration.go ./sql "host=localhost port=5432 user=localhost dbname=db_users sslmode=disable" create add_user_table sql
```


