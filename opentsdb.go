package prom2all


import (
	"github.com/prometheus/prom2json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
	"reflect"
)

// ToOpenTSDBv1 transforms the metrics(not yet Histograms/Summaries) to OpenTSDB line format (v1).
func ToOpenTSDBv1(f *prom2json.Family) string {
	base := fmt.Sprintf("put %s %d", f.Name, time.Now().Unix())
	res := []string{}
	for _, item := range f.Metrics {
		switch item.(type) {
		case prom2json.Metric:
			m := item.(prom2json.Metric)
			val, err := strconv.ParseFloat(m.Value, 64)
			if err != nil {
				continue
			}
			met := fmt.Sprintf("%s %f", base, val)
			if len(m.Labels) != 0 {
				// TODO: Check if key/value meet criteria
				// http://opentsdb.net/docs/build/html/user_guide/writing.html#metrics-and-tags
				lab := []string{}
				for k, v := range m.Labels {
					lab = append(lab, fmt.Sprintf("%s=%s", k, v))
				}
				met = fmt.Sprintf("%s %f %s", base, val, strings.Join(lab, " "))
			}
			res = append(res, met)
		default:
			log.Printf("Type '%s' not yet implemented", reflect.TypeOf(item))
		}
	}
	return strings.Join(res, "\n")
}
