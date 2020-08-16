# pg_flame

A flamegraph generator for Postgres `EXPLAIN ANALYZE` output.

A almost rewrite version,

* static.go is generated with [zipdata](https://github.com/alitrack/zipdata) from static directory.
* Support offline.
* Use [webview](https://github.com/webview/webview) to automatically open created html.
* Support  DATABASE_URL enviroment variable.

<a href="https://mgartner.github.io/pg_flame/flamegraph.html">
  <img width="700" src="https://user-images.githubusercontent.com/1128750/67738754-16f0c300-f9cd-11e9-8fc2-6acc6f288841.png">
</a>

## Demo

Try the demo [here](https://mgartner.github.io/pg_flame/flamegraph.html).

## Installation

### Build from source

If you'd like to build a binary from the source code, run the following
commands. Note that compiling requires Go version 1.13+.

```
$ go get -u github.com/alitrack/pg_flame
```

A `pg_flame` binary will be created that you can place in your `$GOPATH/bin`.

## Usage

```bash
usage: pg_flame [<flags>]

A flamegraph generator for Postgres EXPLAIN ANALYZE output.

Flags:
      --help                    Show context-sensitive help (also try --help-long and --help-man).
  -U, --username="steven"       database user name
  -h, --host="localhost"        database server host or socket directory
  -p, --port=5432               database server port
      --sslmode="disable"       database server sslmode
      --password=PASSWORD       database server password
      --dbname="postgres"       database name
  -o, --output="pg_flame.html"  output html file
  -c, --command=COMMAND         run only single command (SQL)
  -f, --file=FILE               execute commands from file
  -s, --show_browser            Launch browser if successful
      --version                 Show application version.
```

### Example: from sql command 

```bash
$ pg_flame -c  'EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) SELECT id FROM users'
```

### Example: from sql file

Create a SQL file with the `EXPLAIN ANALYZE` query.

```sql
-- query.sql
EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)
SELECT id
FROM users
```

Then run the query and save the JSON to a file.

```bash
$ pg_flame -f query.sql
```

## Background

[Flamegraphs](http://www.brendangregg.com/flamegraphs.html) were invented by
Brendan Gregg to visualize CPU consumption per code-path of profiled software.
They are useful visualization tools in many types of performance
investigations. Flamegraphs have been used to visualize Oracle database
[query
plans](https://blog.tanelpoder.com/posts/visualizing-sql-plan-execution-time-with-flamegraphs/)
and [query
executions](https://externaltable.blogspot.com/2014/05/flame-graphs-for-oracle.html)
, proving useful for debugging slow database queries.

Pg_flame is in extension of that work for Postgres query plans. It generates a
visual hierarchy of query plans. This visualization identifies the relative
time of each part of a query plan.

This tool relies on the
[`spiermar/d3-flame-graph`](https://github.com/spiermar/d3-flame-graph) plugin to
generate the flamegraph.
