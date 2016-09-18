package mo_etcd_client

import (
    "net/http"
    "github.com/prometheus/common/expfmt"
    "etcd_dashboard/modules/mo_conf"
    "time"
    "etcd_dashboard/modules/mo_log"
    //dto "github.com/prometheus/client_model/go"
    "encoding/json"
    "errors"
    "fmt"
)

const monitorDuration = 3 * time.Second

type Monitor struct {
    endpoints []string
}

func NewMonitor() *Monitor {
    endpoints := make([]string, len(mo_conf.Conf.Etcd))

    for i, et := range mo_conf.Conf.Etcd {
        endpoints[i] = et.Addr
    }
    return &Monitor{
        endpoints,
    }
}

func (me *Monitor) StartRock() {
    go func() {
        ticker := time.NewTicker(monitorDuration)
        for range ticker.C {
            me.UpdateMetrics()
        }
    }()
}

func (me *Monitor) getMetrics(d string, mc chan <- []byte) {
    url := "http://" + d + "/metrics"
    resp, err := http.Get(url)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    if err != nil {
        errStr := fmt.Sprint("GET request for URL %q returned HTTP status %s", url, resp.Status)
        err := errors.New(errStr)
        panic(err)
    }

    var parser expfmt.TextParser
    metricFamilies, err2 := parser.TextToMetricFamilies(resp.Body)

    if err2 != nil {
        panic(err2)
    }

    result := []*metricFamily{}
    for _, mf := range metricFamilies {
        //mo_log.Logger.Debug("%v", mf)
        //mc <- mf
        result = append(result, newMetricFamily(mf))
    }
    jsonResult, err := json.Marshal(result)
    mc <- jsonResult
}

func (me *Monitor) UpdateMetrics() {
    //todo: 三组通道
    count := len(me.endpoints)
    mfChan := make(chan []byte, count)
    defer close(mfChan)
    for _, edpt := range me.endpoints {
        go me.getMetrics(edpt, mfChan)
    }
    jsonResults := make([][]byte, count)

    for ; count > 0; count-- {
        tem := <- mfChan
        jsonResults = append(jsonResults, tem)
        mo_log.Logger.Debug(string(tem))
    }

}
