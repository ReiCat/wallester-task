### Commands

Fill the database command:
```sh
$ psql -U {user} -p 5432 -h {ip} -d {database} < sql.sql
```

Run server:
```sh
$ go run main.go
```

Run tests:
```sh
$ go test
```