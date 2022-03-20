package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	dbname        = "homework"
	dbuser        = "postgres"
	maxPoolWorker = 100
)

func main() {

	workerThread := flag.Int("workerThread", 1, fmt.Sprintf("Number of parallel connections to the server %d)", maxPoolWorker))
	file := flag.String("file", "", "CSV file with column format \"hostname, starttime, endtime\" with header intact.")

	flag.Parse()

	if *workerThread < 0 || *workerThread > maxPoolWorker || *file == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	connStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", dbuser, dbname)
	db, err := sql.Open("postgres", connStr)

	checkError(err)
	defer db.Close()

	//validate db connection
	err = db.Ping()
	checkError(err)

	queries := parseCSV(*file + ".csv")

	work := make(chan Request, len(queries))
	defer close(work)

	done := make(chan struct{})
	defer close(done)

	//start workerpool and do the work
	go request(work, queries, done)
	NewScheduler(*workerThread, db).schedule(work, done)
}

func request(work chan Request, queries []QueryParameter, done chan struct{}) {
	c := make(chan Result)

	results := make([]Result, 0, len(queries))

	for _, query := range queries {
		work <- Request{query, c}
		results = append(results, <-c)
	}

	GetStats(results).print()

	done <- struct{}{}
}
