# Prometheus Charts [![Golang CI](https://github.com/ylz-at/prometheus-charts/actions/workflows/go.yml/badge.svg)](https://github.com/ylz-at/prometheus-charts/actions/workflows/go.yml) [![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fylz-at%2Fprometheus-charts.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fylz-at%2Fprometheus-charts?ref=badge_shield)

`prometheus-charts` is used for plotting the metrics in prometheus to images.

### Install
- With Go:
```shell script
$ go get -u github.com/ylz-at/prometheus-charts
```
- From Source:
```shell script
$ go build -o bin/grafana-plot cmd/grafana-plot/main.go
$ go build -o bin/prometheus-plot cmd/prometheus-plot/main.go
```

### Usage

```shell script
$ go run cmd/prometheus-plot/main.go --help
Usage of prometheus-plot:
  -addr string
        Prometheus address (default "http://127.0.0.1:9090")
  -format string
        Output format (default "png")
  -offset value
        Query time offset (default 1h0m0s)
  -output string
        Output file
  -query string
        Query
  -step value
        Query step (default 15s)
  -t value
        Query point time (default Sat Aug  8 17:26:29 CST 2020)
  -title string
        Metric title (default "Metric")
```

### Example

- start prometheus
```shell script
$ docker run --rm -d --name prometheus -p 9090:9090 prometheus/prometheus
```

- run this program
```shell script
$ go run cmd/prometheus-plot/main.go -query go_memstats_frees_total
```

- check the output
```shell script
$ ls Metric.png
Metric.png
```
![Output](Metric.png)



## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fylz-at%2Fprometheus-charts.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fylz-at%2Fprometheus-charts?ref=badge_large)
