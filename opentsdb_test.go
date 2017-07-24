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
	for _,g := range got {
		_, err := regexp.MatchString(p, g)
		assert.NoError(t, err, "Should match OpenTSDBv1 string")
	}
}

func TestSanitizeTags(t *testing.T) {
	got, err := SanitizeTags("tag1", "abcdf")
	assert.NoError(t, err, "fine")
	assert.Equal(t, "tag1=abcdf", got)
	_, err = SanitizeTags("tag2", "abcdf.asda,asd")
	assert.Error(t, err, "fine")
	_, err = SanitizeTags("abcdf.asda,asd", "val3")
	assert.Error(t, err, "fine")
}

func TestLabelToString(t *testing.T) {
	_, err := LabelToString(map[string]string{})
	assert.Error(t, err, "Empty map")
	got, err := LabelToString(map[string]string{"tag1": "val1"})
	assert.NoError(t, err, "fine")
	assert.Equal(t, []string{"tag1=val1"}, got)
	got, err = LabelToString(map[string]string{"tag1": "val1", "tag2": "abcdf.asda,asd"})
	assert.NoError(t, err, "fine")
	assert.Equal(t, []string{"tag1=val1"}, got)
	_, err = LabelToString(map[string]string{"tag2": "abcdf.asda,asd"})
	assert.Error(t, err, "All k/v pairs fail")
}
