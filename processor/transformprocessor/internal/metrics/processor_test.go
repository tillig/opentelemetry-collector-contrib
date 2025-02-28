// Copyright  The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/pmetric"
)

var (
	StartTime      = time.Date(2020, 2, 11, 20, 26, 12, 321, time.UTC)
	StartTimestamp = pcommon.NewTimestampFromTime(StartTime)
	TestTime       = time.Date(2021, 3, 12, 21, 27, 13, 322, time.UTC)
	TestTimeStamp  = pcommon.NewTimestampFromTime(StartTime)
)

func TestProcess(t *testing.T) {
	tests := []struct {
		query []string
		want  func(pmetric.Metrics)
	}{
		{
			query: []string{`set(attributes["test"], "pass") where metric.name == "operationA"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().InsertString("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().InsertString("test", "pass")
			},
		},
		{
			query: []string{`set(attributes["test"], "pass") where resource.attributes["host.name"] == "myhost"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().InsertString("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().InsertString("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().InsertString("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().InsertString("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().InsertString("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().InsertString("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().InsertString("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).Summary().DataPoints().At(0).Attributes().InsertString("test", "pass")

			},
		},
		{
			query: []string{`keep_keys(attributes, "attr2") where metric.name == "operationA"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().InsertString("attr2", "test2")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().InsertString("attr2", "test2")
			},
		},
		{
			query: []string{`set(metric.description, "test") where attributes["attr1"] == "test1"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).SetDescription("test")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).SetDescription("test")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).SetDescription("test")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).SetDescription("test")
			},
		},
		{
			query: []string{`set(metric.unit, "new unit")`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).SetUnit("new unit")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).SetUnit("new unit")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).SetUnit("new unit")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).SetUnit("new unit")
			},
		},
		{
			query: []string{`set(metric.description, "Sum") where metric.type == "Sum"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).SetDescription("Sum")
			},
		},
		{
			query: []string{`set(metric.aggregation_temporality, 1) where metric.aggregation_temporality == 0`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().SetAggregationTemporality(pmetric.MetricAggregationTemporalityDelta)
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().SetAggregationTemporality(pmetric.MetricAggregationTemporalityDelta)
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().SetAggregationTemporality(pmetric.MetricAggregationTemporalityDelta)
			},
		},
		{
			query: []string{`set(metric.is_monotonic, true) where metric.is_monotonic == false`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().SetIsMonotonic(true)
			},
		},
		{
			query: []string{`set(attributes["test"], "pass") where count == 1`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().InsertString("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().InsertString("test", "pass")
			},
		},
		{
			query: []string{`set(attributes["test"], "pass") where scale == 1`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().InsertString("test", "pass")
			},
		},
		{
			query: []string{`set(attributes["test"], "pass") where zero_count == 1`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().InsertString("test", "pass")
			},
		},
		{
			query: []string{`set(attributes["test"], "pass") where positive.offset == 1`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().InsertString("test", "pass")
			},
		},
		{
			query: []string{`set(attributes["test"], "pass") where negative.offset == 1`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().InsertString("test", "pass")
			},
		},
		{
			query: []string{`replace_pattern(attributes["attr1"], "test1", "pass")`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().UpdateString("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().UpdateString("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().UpdateString("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().UpdateString("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().UpdateString("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().UpdateString("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).Summary().DataPoints().At(0).Attributes().UpdateString("attr1", "pass")
			},
		},
		{
			query: []string{`replace_all_patterns(attributes, "test1", "pass")`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().UpdateString("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().UpdateString("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(0).Attributes().UpdateString("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).Histogram().DataPoints().At(1).Attributes().UpdateString("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().UpdateString("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().UpdateString("attr1", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).Summary().DataPoints().At(0).Attributes().UpdateString("attr1", "pass")
			},
		},
		{
			query: []string{`convert_summary_count_val_to_sum("delta", true) where metric.name == "operationD"`},
			want: func(td pmetric.Metrics) {
				sumMetric := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().AppendEmpty()
				sumMetric.SetDataType(pmetric.MetricDataTypeSum)
				sumDp := sumMetric.Sum().DataPoints().AppendEmpty()

				summaryMetric := pmetric.NewMetric()
				fillMetricFour(summaryMetric)
				summaryDp := summaryMetric.Summary().DataPoints().At(0)

				sumMetric.SetDescription(summaryMetric.Description())
				sumMetric.SetName(summaryMetric.Name() + "_count")
				sumMetric.Sum().SetAggregationTemporality(pmetric.MetricAggregationTemporalityDelta)
				sumMetric.Sum().SetIsMonotonic(true)
				sumMetric.SetUnit(summaryMetric.Unit())

				summaryDp.Attributes().CopyTo(sumDp.Attributes())
				sumDp.SetIntVal(int64(summaryDp.Count()))
				sumDp.SetStartTimestamp(StartTimestamp)
				sumDp.SetTimestamp(TestTimeStamp)
			},
		},
		{
			query: []string{`convert_summary_sum_val_to_sum("delta", true) where metric.name == "operationD"`},
			want: func(td pmetric.Metrics) {
				sumMetric := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().AppendEmpty()
				sumMetric.SetDataType(pmetric.MetricDataTypeSum)
				sumDp := sumMetric.Sum().DataPoints().AppendEmpty()

				summaryMetric := pmetric.NewMetric()
				fillMetricFour(summaryMetric)
				summaryDp := summaryMetric.Summary().DataPoints().At(0)

				sumMetric.SetDescription(summaryMetric.Description())
				sumMetric.SetName(summaryMetric.Name() + "_sum")
				sumMetric.Sum().SetAggregationTemporality(pmetric.MetricAggregationTemporalityDelta)
				sumMetric.Sum().SetIsMonotonic(true)
				sumMetric.SetUnit(summaryMetric.Unit())

				summaryDp.Attributes().CopyTo(sumDp.Attributes())
				sumDp.SetDoubleVal(summaryDp.Sum())
				sumDp.SetStartTimestamp(StartTimestamp)
				sumDp.SetTimestamp(TestTimeStamp)
			},
		},
		{
			query: []string{
				`convert_summary_sum_val_to_sum("delta", true) where metric.name == "operationD"`,
				`set(metric.unit, "new unit")`,
			},
			want: func(td pmetric.Metrics) {
				sumMetric := td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().AppendEmpty()
				sumMetric.SetDataType(pmetric.MetricDataTypeSum)
				sumDp := sumMetric.Sum().DataPoints().AppendEmpty()

				summaryMetric := pmetric.NewMetric()
				fillMetricFour(summaryMetric)
				summaryDp := summaryMetric.Summary().DataPoints().At(0)

				sumMetric.SetDescription(summaryMetric.Description())
				sumMetric.SetName(summaryMetric.Name() + "_sum")
				sumMetric.Sum().SetAggregationTemporality(pmetric.MetricAggregationTemporalityDelta)
				sumMetric.Sum().SetIsMonotonic(true)
				sumMetric.SetUnit("new unit")

				summaryDp.Attributes().CopyTo(sumDp.Attributes())
				sumDp.SetDoubleVal(summaryDp.Sum())
				sumDp.SetStartTimestamp(StartTimestamp)
				sumDp.SetTimestamp(TestTimeStamp)

				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).SetUnit("new unit")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(1).SetUnit("new unit")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).SetUnit("new unit")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(3).SetUnit("new unit")
			},
		},
		{
			query: []string{`set(attributes["test"], "pass") where IsMatch(metric.name, "operation[AC]") == true`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().InsertString("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().InsertString("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(0).Attributes().InsertString("test", "pass")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(2).ExponentialHistogram().DataPoints().At(1).Attributes().InsertString("test", "pass")
			},
		},
		{
			query: []string{`delete_key(attributes, "attr3") where metric.name == "operationA"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().InsertString("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().InsertString("attr2", "test2")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().InsertString("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().InsertString("attr2", "test2")
			},
		},
		{
			query: []string{`delete_matching_keys(attributes, "[23]") where metric.name == "operationA"`},
			want: func(td pmetric.Metrics) {
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(0).Attributes().InsertString("attr1", "test1")
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().Clear()
				td.ResourceMetrics().At(0).ScopeMetrics().At(0).Metrics().At(0).Sum().DataPoints().At(1).Attributes().InsertString("attr1", "test1")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.query[0], func(t *testing.T) {
			td := constructMetrics()
			processor, err := NewProcessor(tt.query, DefaultFunctions(), component.ProcessorCreateSettings{})
			assert.NoError(t, err)

			_, err = processor.ProcessMetrics(context.Background(), td)
			assert.NoError(t, err)

			exTd := constructMetrics()
			tt.want(exTd)

			assert.Equal(t, exTd, td)
		})
	}
}

func constructMetrics() pmetric.Metrics {
	td := pmetric.NewMetrics()
	rm0 := td.ResourceMetrics().AppendEmpty()
	rm0.Resource().Attributes().InsertString("host.name", "myhost")
	rm0ils0 := rm0.ScopeMetrics().AppendEmpty()
	fillMetricOne(rm0ils0.Metrics().AppendEmpty())
	fillMetricTwo(rm0ils0.Metrics().AppendEmpty())
	fillMetricThree(rm0ils0.Metrics().AppendEmpty())
	fillMetricFour(rm0ils0.Metrics().AppendEmpty())
	return td
}

func fillMetricOne(m pmetric.Metric) {
	m.SetName("operationA")
	m.SetDescription("operationA description")
	m.SetUnit("operationA unit")
	m.SetDataType(pmetric.MetricDataTypeSum)

	dataPoint0 := m.Sum().DataPoints().AppendEmpty()
	dataPoint0.SetStartTimestamp(StartTimestamp)
	dataPoint0.SetDoubleVal(1.0)
	dataPoint0.Attributes().InsertString("attr1", "test1")
	dataPoint0.Attributes().InsertString("attr2", "test2")
	dataPoint0.Attributes().InsertString("attr3", "test3")

	dataPoint1 := m.Sum().DataPoints().AppendEmpty()
	dataPoint1.SetStartTimestamp(StartTimestamp)
	dataPoint1.Attributes().InsertString("attr1", "test1")
	dataPoint1.Attributes().InsertString("attr2", "test2")
	dataPoint1.Attributes().InsertString("attr3", "test3")
}

func fillMetricTwo(m pmetric.Metric) {
	m.SetName("operationB")
	m.SetDescription("operationB description")
	m.SetUnit("operationB unit")
	m.SetDataType(pmetric.MetricDataTypeHistogram)

	dataPoint0 := m.Histogram().DataPoints().AppendEmpty()
	dataPoint0.SetStartTimestamp(StartTimestamp)
	dataPoint0.Attributes().InsertString("attr1", "test1")
	dataPoint0.Attributes().InsertString("attr2", "test2")
	dataPoint0.Attributes().InsertString("attr3", "test3")
	dataPoint0.SetCount(1)

	dataPoint1 := m.Histogram().DataPoints().AppendEmpty()
	dataPoint1.SetStartTimestamp(StartTimestamp)
	dataPoint1.Attributes().InsertString("attr1", "test1")
	dataPoint1.Attributes().InsertString("attr2", "test2")
	dataPoint1.Attributes().InsertString("attr3", "test3")
}

func fillMetricThree(m pmetric.Metric) {
	m.SetName("operationC")
	m.SetDescription("operationC description")
	m.SetUnit("operationC unit")
	m.SetDataType(pmetric.MetricDataTypeExponentialHistogram)

	dataPoint0 := m.ExponentialHistogram().DataPoints().AppendEmpty()
	dataPoint0.SetStartTimestamp(StartTimestamp)
	dataPoint0.Attributes().InsertString("attr1", "test1")
	dataPoint0.Attributes().InsertString("attr2", "test2")
	dataPoint0.Attributes().InsertString("attr3", "test3")
	dataPoint0.SetCount(1)
	dataPoint0.SetScale(1)
	dataPoint0.SetZeroCount(1)
	dataPoint0.Positive().SetOffset(1)
	dataPoint0.Negative().SetOffset(1)

	dataPoint1 := m.ExponentialHistogram().DataPoints().AppendEmpty()
	dataPoint1.SetStartTimestamp(StartTimestamp)
	dataPoint1.Attributes().InsertString("attr1", "test1")
	dataPoint1.Attributes().InsertString("attr2", "test2")
	dataPoint1.Attributes().InsertString("attr3", "test3")
}

func fillMetricFour(m pmetric.Metric) {
	m.SetName("operationD")
	m.SetDescription("operationD description")
	m.SetUnit("operationD unit")
	m.SetDataType(pmetric.MetricDataTypeSummary)

	dataPoint0 := m.Summary().DataPoints().AppendEmpty()
	dataPoint0.SetStartTimestamp(StartTimestamp)
	dataPoint0.SetTimestamp(TestTimeStamp)
	dataPoint0.Attributes().InsertString("attr1", "test1")
	dataPoint0.Attributes().InsertString("attr2", "test2")
	dataPoint0.Attributes().InsertString("attr3", "test3")
	dataPoint0.SetCount(1234)
	dataPoint0.SetSum(12.34)

	quantileDataPoint0 := dataPoint0.QuantileValues().AppendEmpty()
	quantileDataPoint0.SetQuantile(.99)
	quantileDataPoint0.SetValue(123)

	quantileDataPoint1 := dataPoint0.QuantileValues().AppendEmpty()
	quantileDataPoint1.SetQuantile(.95)
	quantileDataPoint1.SetValue(321)
}
