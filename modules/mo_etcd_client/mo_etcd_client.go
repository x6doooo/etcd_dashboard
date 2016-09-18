package mo_etcd_client

import (
    "github.com/coreos/etcd/clientv3"
    "time"
    "etcd_dashboard/modules/mo_conf"
    "context"
    "reflect"
    "etcd_dashboard/modules/mo_log"
)

const (
    DialTimeout = 10 * time.Second
    RequestTimeout = 10 * time.Second
)

var (
    BaseCli *clientv3.Client
    ClusterCli clientv3.Cluster
    MaintenanceCli clientv3.Maintenance
    Endpoints *EndpointsStruct
    MonitorCli *Monitor
)

func init() {

    endpoints := make([]string, len(mo_conf.Conf.Etcd))

    for i, et := range mo_conf.Conf.Etcd {
        endpoints[i] = et.Addr
    }

    var err error
    BaseCli, err = clientv3.New(clientv3.Config{
        Endpoints: endpoints,
        DialTimeout: DialTimeout,
    })
    if err != nil {
        panic(err)
    }

    ClusterCli = clientv3.NewCluster(BaseCli)
    MaintenanceCli = clientv3.NewMaintenance(BaseCli)

    // 开始循环更新列表数据
    Endpoints = NewEndpointsStuct()
    Endpoints.StartRock()

    st := time.Now()
    ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
    resp, err := BaseCli.Get(ctx, "health")
    cancel()
    mo_log.Logger.Debug("%v, %v, %v", resp, err, time.Since(st))

    // metrics
    //MonitorCli = NewMonitor()
    //MonitorCli.StartRock()
}

func Request(method string, args... interface{}) (resp interface{}, err error) {

    defer func() {
        e := recover()
        if e == nil {
            return
        }
        var ok bool
        err, ok = e.(error)
        if !ok {
            return
        }
    }()

    ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
    defer cancel()

    inputs := make([]reflect.Value, len(args) + 1)
    inputs[0] = reflect.ValueOf(ctx)
    for i, arg := range args {
        inputs[i + 1] = reflect.ValueOf(arg)
    }

    //mo_log.Logger.Debug("%v", ClusterCli)
    //mo_log.Logger.Debug("%v", method)
    //mo_log.Logger.Debug("%v", inputs)
    results := reflect.ValueOf(ClusterCli).MethodByName(method).Call(inputs)

    //mo_log.Logger.Debug("%v", results)
    resp = results[0].Interface()
    errField := results[1].Interface()
    if errField != nil {
        err = errField.(error)
    }

    return
}