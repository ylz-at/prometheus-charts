package plot

import (
	"os"
	"testing"

	"github.com/prometheus/common/model"
)

func TestFile(t *testing.T) {
	type args struct {
		metrics model.Matrix
		title   string
		format  string
		name    string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				metrics: model.Matrix([]*model.SampleStream{
					{
						Metric: model.Metric{"foo": "bar"},
						Values: []model.SamplePair{
							{Timestamp: 0, Value: 0},
							{Timestamp: 1000, Value: 0.5},
							{Timestamp: 2000, Value: 1},
							{Timestamp: 3000, Value: 0.5},
							{Timestamp: 4000, Value: 0},
						},
					},
				}),
				title:  "metrics",
				format: "png",
				name:   "test.png",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := WriteToFile(tt.args.metrics, tt.args.title, tt.args.format, tt.args.name); (err != nil) != tt.wantErr {
				t.Errorf("WriteToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
			info, err := os.Lstat(tt.args.name)
			if err != nil {
				t.Errorf("WriteToFile() error = %v", err)
			}
			if info.Name() != tt.args.name {
				t.Errorf("WriteToFile() error = incorrect name")
			}
		})
	}
}
