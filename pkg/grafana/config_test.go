package grafana

import (
	"reflect"
	"testing"
)

func TestConfigParser(t *testing.T) {
	c, err := ParseConfig("tikv.json")
	if err != nil {
		t.Error(err)
	}

	values := c.LabelValues()

	expectValues := []LabelValue{
		{Name: "db", Query: "tikv_engine_block_cache_size_bytes", Label: "db"},
		{Name: "command", Query: "tikv_storage_command_total", Label: "type"},
	}

	if !reflect.DeepEqual(values, expectValues) {
		t.Errorf("expect label values %v, but got %v", expectValues, values)
	}
}
