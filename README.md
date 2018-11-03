# shutdown.gracefully
This tool with ability  to count summary of jobs with func AddJob and DoneJob. While receive system interrupt  singnal, or application panic to close, this will count down until all jobs done then safety return.

Check [example](https://github.com/czminami/shutdown.gracefully/blob/master/example/shutdown_example.go) for how to use.