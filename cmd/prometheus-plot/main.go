package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/prometheus/common/model"

	"github.com/ylz-at/prometheus-charts/pkg/flags"
	"github.com/ylz-at/prometheus-charts/pkg/plot"
	"github.com/ylz-at/prometheus-charts/pkg/prometheus"
)

var (
	addr      = flag.String("addr", "http://127.0.0.1:9090", "Prometheus address")
	query     = flag.String("query", "", "Query")
	pointTime = flags.UnixTime("t", time.Now(), "Query point time")
	offset    = flags.Duration("offset", time.Hour, "Query time offset")
	step      = flags.Duration("step", 15*time.Second, "Query step")
	output    = flag.String("output", "", "Output file")
	format    = flag.String("format", "png", "Output format")
	title     = flag.String("title", "Metric", "Metric title")
)

func main() {
	flag.Parse()

	c, err := prometheus.NewClient(*addr)
	if err != nil {
		panic(err)
	}

	startTime, endTime := (*pointTime).Add(-*offset), *pointTime

	v, _, err := c.Query(context.Background(), *query, startTime, endTime, *step)
	if err != nil {
		panic(err)
	}

	m, ok := v.(model.Matrix)
	if !ok {
		panic("must support matrix metric")
	}

	if len(*output) == 0 {
		*output = fmt.Sprintf("%s.%s", *title, *format)
	}

	if err := plot.PlotFile(m, *title, *format, *output); err != nil {
		panic(err)
	}
}
