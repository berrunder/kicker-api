# kicker-api

Simple REST API for database of kicker (foosball) games

### Installation
1. Install [Go](https://golang.org/doc/install);
1. (Optional) Run docker  container for mysql database with `docker-compose up`;
1. Set database credentials for your environment in `dbconfig.yaml`;
1. Install [sql-migrate](https://github.com/rubenv/sql-migrate) tool:

  ```bash
  go get github.com/rubenv/sql-migrate/...
  ```
1. Run migrations on database (for production use `-env=production` flag):

  ```bash
  sql-migrate up
  ```
6. Build with `go build` in kicker-api directory.

### Usage
```
kicker-api [--port value] [--datasource value] [--dialect value]
```
Where `--port` is webserver port (4000 by default), `--datasource` is database connection string (by default DSN from dbconfig.yaml is used), `--dialect` is database dialect ("mysql" by default, can be any dialect supported by sqlx).
`KICKER_ENV` environment variable is used for choosing database config from dbconfig.yaml ("development" by default). 
