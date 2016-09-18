package mo_etcd_client

import (
    "fmt"
    dto "github.com/prometheus/client_model/go"
)

type metricFamily struct {
    Name    string        `json:"name"`
    Help    string        `json:"help"`
    Type    string        `json:"type"`
    Metrics []interface{} `json:"metrics,omitempty"` // Either metric or summary.
}

// metric is for all "single value" metrics.
type metric struct {
    Labels map[string]string `json:"labels,omitempty"`
    Value  string            `json:"value"`
}

type summary struct {
    Labels    map[string]string `json:"labels,omitempty"`
    Quantiles map[string]string `json:"quantiles,omitempty"`
    Count     string            `json:"count"`
    Sum       string            `json:"sum"`
}

type histogram struct {
    Labels  map[string]string `json:"labels,omitempty"`
    Buckets map[string]string `json:"buckets,omitempty"`
    Count   string            `json:"count"`
    Sum     string            `json:"sum"`
}

func newMetricFamily(dtoMF *dto.MetricFamily) *metricFamily {
    mf := &metricFamily{
        Name:    dtoMF.GetName(),
        Help:    dtoMF.GetHelp(),
        Type:    dtoMF.GetType().String(),
        Metrics: make([]interface{}, len(dtoMF.Metric)),
    }
    for i, m := range dtoMF.Metric {
        if dtoMF.GetType() == dto.MetricType_SUMMARY {
            mf.Metrics[i] = summary{
                Labels:    makeLabels(m),
                Quantiles: makeQuantiles(m),
                Count:     fmt.Sprint(m.GetSummary().GetSampleCount()),
                Sum:       fmt.Sprint(m.GetSummary().GetSampleSum()),
            }
        } else if dtoMF.GetType() == dto.MetricType_HISTOGRAM {
            mf.Metrics[i] = histogram{
                Labels:  makeLabels(m),
                Buckets: makeBuckets(m),
                Count:   fmt.Sprint(m.GetHistogram().GetSampleCount()),
                Sum:     fmt.Sprint(m.GetSummary().GetSampleSum()),
            }
        } else {
            mf.Metrics[i] = metric{
                Labels: makeLabels(m),
                Value:  fmt.Sprint(getValue(m)),
            }
        }
    }
    return mf
}

func getValue(m *dto.Metric) float64 {
    if m.Gauge != nil {
        return m.GetGauge().GetValue()
    }
    if m.Counter != nil {
        return m.GetCounter().GetValue()
    }
    if m.Untyped != nil {
        return m.GetUntyped().GetValue()
    }
    return 0.
}

func makeLabels(m *dto.Metric) map[string]string {
    result := map[string]string{}
    for _, lp := range m.Label {
        result[lp.GetName()] = lp.GetValue()
    }
    return result
}

func makeQuantiles(m *dto.Metric) map[string]string {
    result := map[string]string{}
    for _, q := range m.GetSummary().Quantile {
        result[fmt.Sprint(q.GetQuantile())] = fmt.Sprint(q.GetValue())
    }
    return result
}

func makeBuckets(m *dto.Metric) map[string]string {
    result := map[string]string{}
    for _, b := range m.GetHistogram().Bucket {
        result[fmt.Sprint(b.GetUpperBound())] = fmt.Sprint(b.GetCumulativeCount())
    }
    return result
}

