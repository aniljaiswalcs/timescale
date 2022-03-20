package main

import (
	"testing"
	"time"
)

type testParameter struct {
	host      string
	starttime time.Time
	endtime   time.Time
}

type testParsecsv struct {
	Items []testParameter
}

func (q *testParsecsv) generateTestCsvData() {
	hname := make([]string, 0)
	hname = append(hname, "host_000003", "host_000005", "")
	starttime := make([]string, 0)
	starttime = append(starttime, "2017-01-01 08:52:14", "2017-01-01 21:45:18", "2017-01-01 21:45:18")
	endtime := make([]string, 0)
	endtime = append(endtime, "2017-01-01 09:52:14", "2017-01-01 22:45:18", "2017-01-01 22:45:18")

	for i := 0; i < len(hname); i++ {
		t1, err := time.Parse(timeStandard, starttime[i])
		if err != nil {
			panic(err)
		}
		t2, errr := time.Parse(timeStandard, endtime[i])
		if errr != nil {
			panic(errr)
		}
		q.Items = append(q.Items, testParameter{host: hname[i], starttime: t1, endtime: t2})
	}

}

func TestParseCSV(t *testing.T) {
	r1 := make([]QueryParameter, 10)

	sl := testParsecsv{}
	sl.generateTestCsvData()

	tests := []struct {
		testname       string
		expectedResult testParameter
	}{
		{testname: "host_000003,2017-01-01 08:52:14,2017-01-01 09:52:14", expectedResult: testParameter{host: sl.Items[0].host, starttime: sl.Items[0].starttime, endtime: sl.Items[0].endtime}},
		{testname: "host_000005,2017-01-01 21:45:18,2017-01-01 22:45:18", expectedResult: testParameter{host: sl.Items[1].host, starttime: sl.Items[1].starttime, endtime: sl.Items[1].endtime}},
		{testname: " host name blank: ' ' ,2017-01-01 21:45:18,2017-01-01 22:45:18", expectedResult: testParameter{host: sl.Items[2].host, starttime: sl.Items[2].starttime, endtime: sl.Items[2].endtime}},
	}
	for index, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			file := "csv_test.csv"
			r1 = parseCSV(file)
			if (r1[index].Hostname != tt.expectedResult.host) || (r1[index].Starttime != tt.expectedResult.starttime) || (r1[index].Endtime != tt.expectedResult.endtime) {
				t.Errorf("Expected Hostname = %v, Got the Hostname =  %v and Expected starttime = %v , and Got the starttime = %v Expected endtime = %v ,Got the endtime = %v \n", r1[index].Hostname, tt.expectedResult.host, r1[index].Starttime, tt.expectedResult.starttime, r1[index].Endtime, tt.expectedResult.endtime)
			}
		})
	}

}
