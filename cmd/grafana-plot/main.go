package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/prometheus/common/model"
	"github.com/ylz-at/prometheus-charts/pkg/plot"

	"github.com/ylz-at/prometheus-charts/pkg/prometheus"

	"github.com/ylz-at/prometheus-charts/pkg/flags"
	"github.com/ylz-at/prometheus-charts/pkg/grafana"
)

var (
	addr      = flag.String("addr", "http://127.0.0.1:9090", "Prometheus address")
	pointTime = flags.UnixTime("t", time.Now(), "Query time")
	offset    = flags.Duration("offset", time.Hour, "Query time offset, the time range is [time - offset, time]")
	step      = flags.Duration("step", 15*time.Second, "Query step")
	output    = flag.String("output", "./var", "Output dir")
	format    = flag.String("format", "png", "Output format")
	jsonFile  = flag.String("f", "", "Grafana JSON config file")
)

func perr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	os.RemoveAll(*output)

	err := os.MkdirAll(*output, 0755)
	perr(err)

	startTime, endTime := (*pointTime).Add(-*offset), *pointTime

	cfg, err := grafana.ParseConfig(*jsonFile)
	perr(err)

	client, err := prometheus.NewClient(*addr)
	perr(err)

	labelValues := map[string]model.LabelValues{}

	ctx := context.Background()

	for _, v := range cfg.LabelValues() {
		values, err := client.LabelValues(ctx, v.Query, v.Label, startTime, endTime)
		perr(err)
		labelValues[v.Name] = values
	}

	fileID := 0
	for _, row := range cfg.Rows {
		rows := getRepeatRows(row.Title, row.Repeat, labelValues)
		if rows == nil {
			rows = model.LabelValues{""}
		}

		for _, r := range rows {
			rowTitle := replaceQueryLabelValues(row.Title, row.Repeat, string(r))

			for _, panel := range row.Panels {
				fileID++
				name := fmt.Sprintf("%s/%04d-%s-%s.%s", *output, fileID, rowTitle, panel.Title, *format)

				var matrix model.Matrix
				for _, expr := range panel.Targets {
					query := replaceQueryLabelValues(expr.Expr, row.Repeat, string(r))
					v, _, err := client.Query(ctx, query, startTime, endTime, *step)
					if err != nil {
						perr(fmt.Errorf("query %s failed %v", query, err))
					}

					m, ok := v.(model.Matrix)
					if !ok {
						fmt.Printf("query %s return not matrix, skip\n", query)
						continue
					}
					matrix = append(matrix, m...)
				}

				err = plot.WriteToFile(matrix, panel.Title, *format, name)
				if err != nil {
					// TODO: handle Nan data point
					fmt.Printf("plot %s-%s failed %v\n", rowTitle, panel.Title, err)
				}
			}
		}
	}
}

func getRepeatRows(title string, repeat string, labelValues map[string]model.LabelValues) model.LabelValues {
	if len(repeat) == 0 || !strings.Contains(title, "$"+repeat) {
		return nil
	}

	return labelValues[repeat]
}

func replaceQueryLabelValues(query string, label string, value string) string {
	if len(label) == 0 || len(value) == 0 {
		return query
	}

	re := regexp.MustCompile(`\$` + label)
	return re.ReplaceAllString(query, value)
}
