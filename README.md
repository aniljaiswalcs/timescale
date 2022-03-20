# Benchmarking

select query benchmarking for TimeScaleDB

Benchmarking takes the number of workers and a .csv file with format "hostname, start_time, end_time" (header included) to generate and benchmark query at run time. Queries are allotted to workers based on hostnames (no two workers share queries that touch the same hostnames).

## Build and Usage

First set up PostgreSql server.
run timescale on top of PostgreSql
insert data into timescale 

```bash
clone the repo

go build -o selectbenchmark  

./selectbenchmark -workerThread 10 -file query_params

go test -v
```
Testing:
Few testcase added.
