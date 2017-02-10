# pgsplit
Split a postgresql script file into statements

`pgsplit` takes a PostgreSQL format sql script on stdin, and spits it out
on stdout with each statement separated by a zero byte.

This is for database migration code that wants to send each statement to the
database separately, rather than sending the entire script as a giant
semicolon separated string.

It could probably be abused along with `xargs -0` to do something interesting
too.

It's very simple. It splits the file at each line that ends with a semicolon
that isn't inside a quoted string.

# What works

Strings quoted with single quotes or [dollar-quoting](https://www.postgresql.org/docs/current/static/sql-syntax-lexical.html#SQL-SYNTAX-DOLLAR-QUOTING).

Nested quotes.

My SQL schema upgrade scripts.

# What doesn't

COPY. Backslash-escaped single quotes.

Your SQL schema upgrade scripts (probably).
