package prometheus

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// Client queries the prometheus metrics
type Client struct {
	api v1.API
}

// NewClient creates the prometheus client.
func NewClient(addr string) (*Client, error) {
	c, err := api.NewClient(api.Config{Address: addr})
	if err != nil {
		return nil, err
	}
	return &Client{api: v1.NewAPI(c)}, nil
}

// Query queries the metric from start to end time with step.
func (c *Client) Query(ctx context.Context, query string, startTime time.Time, endTime time.Time, step time.Duration) (model.Value, v1.Warnings, error) {
	return c.api.QueryRange(ctx, query, v1.Range{Start: startTime, End: endTime, Step: step})
}

// Series finds series by label matchers.
func (c *Client) Series(ctx context.Context, matches []string, startTime time.Time, endTime time.Time) ([]model.LabelSet, v1.Warnings, error) {
	return c.api.Series(ctx, matches, startTime, endTime)
}

// LabelValues returns all the values for the label by the matchers
func (c *Client) LabelValues(ctx context.Context, matcher string, label string, startTime time.Time, endTime time.Time) (model.LabelValues, error) {
	set, _, err := c.Series(ctx, []string{matcher}, startTime, endTime)
	if err != nil {
		return nil, err
	}

	lvs := make(model.LabelValues, 0, len(set))
	for _, s := range set {
		if v, ok := s[model.LabelName(label)]; ok {
			lvs = append(lvs, v)
		}
	}

	return lvs, nil
}
