# Scribo Database Migrations

The database migrations scheme for Scribo is very simple. Rather than implement a whole hoopla with creating a management scheme, I've written a go command called scribo-migrate which simply looks in this directory for the latest migration, as defined by the filename, then runs the SQL in that script.

This means that the database does not account for its own version, nor is there any upgrade or downgrade command. Like I said, very simple. If multiple migrations need to be run, they can be specified as arguments to the script.

Consider a migrations directory containing the following files:

```
- 0001-initialize.sql
- 0002-create-node-table.sql
- 0003-create-ping-table.sql
- 0004-alter-node-data.sql
```

Running `scribo-migrate` would execute the SQL in the `0004-alter-node-data.sql` file. This is equivalent to executing `scribo-migrate --latest`. Alternatively you can specify the file to migrate: `scribo-migrate 1` would run `0001-initialize.sql`. Note that the command will only identify SQL files that begin with a number, and will access them by number.
