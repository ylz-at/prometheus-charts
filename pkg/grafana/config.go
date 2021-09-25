package grafana

import (
	"encoding/json"
	"io/ioutil"
	"strings"
)

// PanelConfig is for panel of the row.
type PanelConfig struct {
	Targets []struct {
		Expr         string `json:"expr"`
		LegendFormat string `json:"legendFormat"`
	}
	Title string `json:"title"`
}

// Config is the Prometheus metrics config for the Grafana,
// We will parse the associated json file to this Config.
type Config struct {
	Rows []struct {
		Panels []PanelConfig `json:"panels"`
		Repeat string        `json:"repeat"`
		Title  string        `json:"title"`
	} `json:"rows"`
	Templating struct {
		List []struct {
			Name  string `json:"name"`
			Label string `json:"label"`
			Query string `json:"query"`
		} `json:"list"`
	} `json:"templating"`
}

// ParseConfig parses the config from a JSON file.
func ParseConfig(path string) (*Config, error) {
	var c Config

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal(data, &c); err != nil {
		return nil, err
	}

	return &c, nil
}

// LabelValue is used to get thet time series with the centain label.
// We will use this label value to replace the name configured in Grafana.
type LabelValue struct {
	Query string
	Label string
	Name  string
}

// LabelValues parse the templating to get the label values.
func (c *Config) LabelValues() []LabelValue {
	values := make([]LabelValue, 0, len(c.Templating.List))

	for _, v := range c.Templating.List {
		if !strings.HasPrefix(v.Query, "label_values(") {
			continue
		}

		q := strings.Trim(v.Query, "label_values(")
		q = strings.Trim(q, ")")
		seps := strings.Split(q, ", ")
		if len(seps) != 2 {
			continue
		}

		values = append(values, LabelValue{
			Name:  v.Name,
			Query: seps[0],
			Label: seps[1],
		})
	}

	return values
}
