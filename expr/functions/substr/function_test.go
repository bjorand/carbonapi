package substr

import (
	"testing"
	"time"

	"github.com/go-graphite/carbonapi/expr/helper"
	"github.com/go-graphite/carbonapi/expr/metadata"
			th "github.com/go-graphite/carbonapi/tests"
	"github.com/go-graphite/carbonapi/pkg/parser"
	"github.com/go-graphite/carbonapi/expr/types"
)

func init() {
	md := New("")
	evaluator := th.EvaluatorFromFunc(md[0].F)
	metadata.SetEvaluator(evaluator)
	helper.SetEvaluator(evaluator)
	for _, m := range md {
		metadata.RegisterFunction(m.Name, m.F)
	}
}

func TestSubstr(t *testing.T) {
	now32 := int64(time.Now().Unix())

	/*
	Python's behavior:
	>>> a = ["metric1", "foo", "bar", "baz"]
	>>> a[-3:-1]
	['foo', 'bar']
	>>> a[-4:-1]
	['metric1', 'foo', 'bar']
	>>> a[-65:]
	['metric1', 'foo', 'bar', 'baz']
	>>> a[-6:-1]
	['metric1', 'foo', 'bar']
	>>> a[0:-1]
	['metric1', 'foo', 'bar']
	 */
	tests := []th.EvalTestItem{
		{
			parser.NewExpr("substr",
				"metric1.foo.bar.baz", 1, 3,
			),
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1.foo.bar.baz", 0, 1}: {types.MakeMetricData("metric1.foo.bar.baz", []float64{1, 2, 3, 4, 5}, 1, now32)},
			},
			[]*types.MetricData{types.MakeMetricData("foo.bar",
				[]float64{1, 2, 3, 4, 5}, 1, now32)},
		},
		{
			parser.NewExpr("substr",
				"metric1.foo.bar.baz", -3, -1,
			),
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1.foo.bar.baz", 0, 1}: {types.MakeMetricData("metric1.foo.bar.baz", []float64{1, 2, 3, 4, 5}, 1, now32)},
			},
			[]*types.MetricData{types.MakeMetricData("foo.bar",
				[]float64{1, 2, 3, 4, 5}, 1, now32)},
		},
		{
			parser.NewExpr("substr",
				"metric1.foo.bar.baz", -3,
			),
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1.foo.bar.baz", 0, 1}: {types.MakeMetricData("metric1.foo.bar.baz", []float64{1, 2, 3, 4, 5}, 1, now32)},
			},
			[]*types.MetricData{types.MakeMetricData("foo.bar.baz",
				[]float64{1, 2, 3, 4, 5}, 1, now32)},
		},
		{
			parser.NewExpr("substr",
				"metric1.foo.bar.baz", -6, -1,
			),
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1.foo.bar.baz", 0, 1}: {types.MakeMetricData("metric1.foo.bar.baz", []float64{1, 2, 3, 4, 5}, 1, now32)},
			},
			[]*types.MetricData{types.MakeMetricData("metric1.foo.bar",
				[]float64{1, 2, 3, 4, 5}, 1, now32)},
		},
		{
			parser.NewExpr("substr",
				"metric1.foo.bar.baz", 0, -1,
			),
			map[parser.MetricRequest][]*types.MetricData{
				{"metric1.foo.bar.baz", 0, 1}: {types.MakeMetricData("metric1.foo.bar.baz", []float64{1, 2, 3, 4, 5}, 1, now32)},
			},
			[]*types.MetricData{types.MakeMetricData("metric1.foo.bar",
				[]float64{1, 2, 3, 4, 5}, 1, now32)},
		},
	}

	for _, tt := range tests {
		testName := tt.E.Target() + "(" + tt.E.RawArgs() + ")"
		t.Run(testName, func(t *testing.T) {
			th.TestEvalExpr(t, &tt)
		})
	}

}
