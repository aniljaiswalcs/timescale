package main

import (
	"fmt"
	"sort"
	"time"
)

// Stats encapsulates the stats for a list of time.Duration values
type StattoPrint struct {
	Count   int
	Total   ResultTiming
	Average float64
	Median  float64
	Minimum ResultTiming
	Maximum ResultTiming
}

// ResultTiming represents the time spent executing
type ResultTiming struct {
	time    float64
	results int
}

func (rt *ResultTiming) stringvalue() string {
	return fmt.Sprintf("%.4f ms (%d results)", rt.time, rt.results)
}

func (s *StattoPrint) print() {

	fmt.Println("")
	fmt.Printf("Total No. of queries:        %d \n", s.Count)
	fmt.Printf("Total time taken to execute: %s\n", s.Total.stringvalue())
	fmt.Printf("Average time execution:      %.4f ms\n", s.Average)
	fmt.Printf("Median time execution:       %.4f ms\n", s.Median)
	fmt.Printf("Minimum time execution:      %s\n", s.Minimum.stringvalue())
	fmt.Printf("Maximum time execution:      %s\n", s.Maximum.stringvalue())
	fmt.Println("")

}

//creates and returns the Stat
func GetStats(results []Result) *StattoPrint {
	stat := &StattoPrint{}
	var n int

	durationToNS := func(d time.Duration) float64 {
		return float64(d.Nanoseconds()) / float64(1000000)
	}

	resTimes := Map(results, durationToNS)

	sort.Slice(resTimes, func(i, j int) bool {
		return resTimes[i].time < resTimes[j].time
	})

	n = len(resTimes)

	if n == 0 {
		return stat
	}

	stat.Count = n
	stat.Minimum = resTimes[0]
	stat.Maximum = resTimes[n-1]

	totalTime := 0.0
	totalRes := 0

	for _, rt := range resTimes {
		totalTime += rt.time
		totalRes += rt.results
	}

	stat.Total = ResultTiming{totalTime, totalRes}

	stat.Average = totalTime / float64(len(resTimes))

	if n%2 == 0 && n > 1 {
		stat.Median = (resTimes[n/2].time + resTimes[n/2+1].time) / 2
	} else {
		stat.Median = resTimes[n/2].time
	}

	return stat
}
