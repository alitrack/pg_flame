package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"strings"

	"github.com/alecthomas/kingpin"
	"github.com/alitrack/pg_flame/pkg/html"
	"github.com/alitrack/pg_flame/pkg/plan"
	"github.com/jackc/pgx/v4"
)

var (
	usr, _      = user.Current()
	userName    = kingpin.Flag("username", "database user name").Default(usr.Username).Short('U').String()
	host        = kingpin.Flag("host", "database server host or socket directory").Short('h').Default("localhost").String()
	port        = kingpin.Flag("port", "database server port").Short('p').Default("5432").Int()
	sslmode     = kingpin.Flag("sslmode", "database server sslmode").Default("disable").String()
	password    = kingpin.Flag("password", "database server password").String()
	dbName      = kingpin.Flag("dbname", "database name").Default("postgres").String()
	output      = kingpin.Flag("output", "output html file").Default("pg_flame.html").Short('o').String()
	sql         = kingpin.Flag("command", "run only single command (SQL)").Short('c').String()
	file        = kingpin.Flag("file", "execute commands from file").Short('f').ExistingFile()
	showBrowser = kingpin.Flag("show_browser", "Launch browser if successful").Default("true").Short('s').Bool()
	version     = "0.1.0"
)

func main() {
	kingpin.Version(version)

	kingpin.CommandLine.Help = "A flamegraph generator for Postgres EXPLAIN ANALYZE output."
	kingpin.Parse()

	dbURL := fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s user=%s password=%s",
		*host, *port, *dbName, *sslmode, *userName, *password)

	conn, err := pgx.Connect(context.Background(), dbURL)
	if err != nil {
		conn, err = pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
			kingpin.Usage()
			return
		}
	}
	defer conn.Close(context.Background())

	// sql := "EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON) SELECT * FROM iris"

	switch {
	case *sql != "":
		break
	case *file != "":
		bs, err := ioutil.ReadFile(*file)
		if err != nil {
			fmt.Println(err)
			return
		}
		*sql = string(bs)
	default:
		fmt.Println("Please use sql command(-c) or sql file(-f) ")
		return
	}

	var json string
	err = conn.QueryRow(context.Background(), *sql).Scan(&json)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}

	reader := strings.NewReader(json)

	p, err := plan.New(reader)
	if err != nil {
		fmt.Println(err)
	}

	o := os.Stdout
	if *output != "" {
		o, err = os.Create(*output)
		if err != nil {
			fmt.Println(err)
		}
	}
	err = html.Generate(o, p)
	if err != nil {
		fmt.Println(err)
	}
	if *showBrowser && *output != "" {
		view(*output)
	}
}
