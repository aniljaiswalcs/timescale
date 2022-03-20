package main

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	_ "github.com/lib/pq"
)

const (
	namedb       = "homework"
	userdb       = "postgres"
	Workermax    = 100
	timeStandard = "2006-01-05 15:04:05"
)

type testRequest struct {
	testquery  testQuery
	testresult chan testResult
}

// Query defines the model for building the queries
type testQuery struct {
	testHostname  string
	testStarttime time.Time
	testEndtime   time.Time
}

type testResult struct {
	testcount    int
	testduration time.Duration
}
type slTestQuery struct {
	Items []testQuery
}

func (q *slTestQuery) generatetestquery() {
	hname := make([]string, 0)
	hname = append(hname, "host_000010", "host_000007", "host_000001", "host_000000", "host_000006", "host_000008", "host_000005", "host_000015")
	starttime := make([]string, 0)
	starttime = append(starttime, "2017-01-01 08:59:22", "2017-01-02 09:15:09", "2017-01-01 08:59:22", "2017-01-02 15:44:45", "2017-01-01 04:30:52", "2017-02-01 18:50:28", "2020-01-02 02:28:31", "2023-01-02 02:28:31")
	endtime := make([]string, 0)
	endtime = append(endtime, "2017-01-01 09:01:09", "2017-01-02 10:15:09", "2017-01-01 09:01:09", "2017-01-02 16:44:45", "2017-01-01 05:30:52", "2017-02-01 19:50:28", "2020-01-02 03:28:31", "2023-01-02 20:28:31")

	for i := 0; i < len(hname); i++ {
		t1, err := time.Parse(timeStandard, starttime[i])
		if err != nil {
			panic(err)
		}
		t2, errr := time.Parse(timeStandard, endtime[i])
		if errr != nil {
			panic(errr)
		}
		q.Items = append(q.Items, testQuery{testHostname: hname[i], testStarttime: t1, testEndtime: t2})
	}

}

const (
	testbaseQuery = `SELECT time_bucket('1 minute', ts) as minute,
	MAX(usage) as max_usage,
	MIN(usage) as min_usage
	FROM cpu_usage
	WHERE host = $1 AND ts BETWEEN $2 AND $3
	GROUP BY minute`
)

func TestRunQuery(t *testing.T) {

	conStr := fmt.Sprintf("user=%s dbname=%s sslmode=disable", userdb, namedb)
	dbstr, _ := sql.Open("postgres", conStr)
	parsed, err := dbstr.Prepare(testbaseQuery)
	checkError(err)

	defer dbstr.Close()
	sl := slTestQuery{}
	sl.generatetestquery()

	tests := []struct {
		testname       string
		inputQuery     testQuery
		expectedResult testResult
	}{
		{testname: "3 rows between 2017-01-01 08:59:22 and 2017-01-01 09:01:09 over 1 min", inputQuery: testQuery{testHostname: sl.Items[0].testHostname, testStarttime: sl.Items[0].testStarttime, testEndtime: sl.Items[0].testEndtime}, expectedResult: testResult{testcount: 3, testduration: 3}},
		{testname: "61 rows between 017-01-01 09:15:09 and 2017-01-01 10:15:09  over 1 min", inputQuery: testQuery{testHostname: sl.Items[1].testHostname, testStarttime: sl.Items[1].testStarttime, testEndtime: sl.Items[1].testEndtime}, expectedResult: testResult{testcount: 61, testduration: 3}},
		{testname: "3 rows between 2017-01-01 08:59:22 and 2017-01-01 09:01:09 over 1 min", inputQuery: testQuery{testHostname: sl.Items[2].testHostname, testStarttime: sl.Items[2].testStarttime, testEndtime: sl.Items[2].testEndtime}, expectedResult: testResult{testcount: 3, testduration: 3}},
		{testname: "61 rows for host_000000 between 2017-01-01 15:44:45 and 2017-01-01 16:44:45 over 1 min", inputQuery: testQuery{testHostname: sl.Items[3].testHostname, testStarttime: sl.Items[3].testStarttime, testEndtime: sl.Items[3].testEndtime}, expectedResult: testResult{testcount: 61, testduration: 3}},
		{testname: "60 rows for host_000006 between 2017-01-01 04:30:52 and 2017-01-01 05:30:52 over 1 min", inputQuery: testQuery{testHostname: sl.Items[4].testHostname, testStarttime: sl.Items[4].testStarttime, testEndtime: sl.Items[4].testEndtime}, expectedResult: testResult{testcount: 60, testduration: 3}},
		{testname: "0 rows for host_000008 between 2017-01-02 18:50:28 and 2017-01-02 19:50:28 over 1 min", inputQuery: testQuery{testHostname: sl.Items[5].testHostname, testStarttime: sl.Items[5].testStarttime, testEndtime: sl.Items[5].testEndtime}, expectedResult: testResult{testcount: 0, testduration: 3}},
		{testname: "0 rows for host_000008 between 2020-01-02 02:28:31 and 2017-01-02 19:50:28 over 1 min", inputQuery: testQuery{testHostname: sl.Items[6].testHostname, testStarttime: sl.Items[6].testStarttime, testEndtime: sl.Items[6].testEndtime}, expectedResult: testResult{testcount: 0, testduration: 3}},
		{testname: "0 rows for host_000015 between 2023-01-02 02:28:31 and 2023-01-02 22:28:31 over 1 min", inputQuery: testQuery{testHostname: sl.Items[7].testHostname, testStarttime: sl.Items[7].testStarttime, testEndtime: sl.Items[7].testEndtime}, expectedResult: testResult{testcount: 0, testduration: 3}},
	}
	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			resultCount := 0
			start := time.Now()
			var duration time.Duration
			rows, err := parsed.Query(tt.inputQuery.testHostname, tt.inputQuery.testStarttime, tt.inputQuery.testEndtime)
			if err == nil {
				for rows.Next() {
					duration = time.Now().Sub(start)
					resultCount++
				}
			}
			if resultCount != tt.expectedResult.testcount && duration != tt.expectedResult.testduration {
				t.Errorf("output count = %v, expecteds %v and duration = %v , expected duration = %v", resultCount, tt.expectedResult.testcount, duration, tt.expectedResult.testcount)
			}
		})
	}
}
