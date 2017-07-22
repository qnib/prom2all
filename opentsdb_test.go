package prom2all

import (
	"regexp"
	"testing"
	"github.com/stretchr/testify/assert"
	dto "github.com/prometheus/client_model/go"

	"github.com/prometheus/prom2json"
)
func strPtr(s string) *string {
	return &s
}
func floatPtr(f float64) *float64 {
	return &f
}
func metricTypePtr(mt dto.MetricType) *dto.MetricType {
	return &mt
}
func createLabelPair(name string, value string) *dto.LabelPair {
	return &dto.LabelPair{
		Name:  &name,
		Value: &value,
	}
}
/// Testvalues
var (
	// Test metric with tags
	mWt1 =  &dto.Metric{
		Label: []*dto.LabelPair{
			createLabelPair("tag1", "abc"),
			createLabelPair("tag2", "def"),
		},
		Counter: &dto.Counter{
			Value: floatPtr(1),
		},
	}
	// Test metric without tags
	mWOt1 = &dto.Metric{
		Label: []*dto.LabelPair{},
		Counter: &dto.Counter{
			Value: floatPtr(2),
		},
	}
  	mf1 =   &dto.MetricFamily{
		Name: strPtr("counter1"),
		Type: metricTypePtr(dto.MetricType_COUNTER),
		Metric: []*dto.Metric{
			mWt1,
			mWOt1,
		},
	}
)

func TestFamily_ToOpenTSDBv1(t *testing.T) {
	x1 := prom2json.NewFamily(mf1)
	got := ToOpenTSDBv1(x1)
	p := `put counter1 [0-9]+ 1.000000 tag1=abc tag2=def\nput counter1 [0-9]+ 2.000000`
	_, err := regexp.MatchString(p, got)
	assert.NoError(t, err, "Should match OpenTSDBv1 string")
}
