# Backend

Create a `./config.toml` file based on the `./config.example.toml` one.  
You can use a Sqlite3 or MySQL database, check the example config.  

For sqlite, you can use the `schema.sqlite.sql` file to create the tables
and the `seed.sql` file to insert some example data:

```sh
sqlite3 ./testdb < schema.sqlite.sql
sqlite3 ./testdb < seed.sql
```

If you use MySQL, leave `host = "db"` to use the docker container or replace
it with some address (like `127.0.0.1`) to use another mysql db.  
To create the tables you can use the `schema.mysql.sql` file:

```sh
echo "CREATE DATABASE fattingo_test;" | mysql -uroot
mysql -uroot fattingo_test < schema.mysql.sql
mysql -uroot fattingo_test < seed.sql
```
