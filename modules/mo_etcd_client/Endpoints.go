package mo_etcd_client

import (
    "github.com/coreos/etcd/clientv3"
    "context"
    "sync"
    "etcd_dashboard/modules/mo_log"
    "time"
    "github.com/coreos/etcd/etcdserver/etcdserverpb"
    "etcd_dashboard/modules/mo_conf"
)

const (
    tickDuration = 10 * time.Second
)

type EndpointDetail struct {
    Member *etcdserverpb.Member
    AddrForDashboard string
    Version string
    DbSize int64
    Leader uint64
    RaftIndex uint64
    RaftTerm uint64
    Latency string
    IsHealth bool
}

type EndpointsStruct struct {
    List []EndpointDetail
    mu sync.Mutex
}

func NewEndpointsStuct() *EndpointsStruct {
    return &EndpointsStruct{
        List: make([]EndpointDetail, 0),
        mu: sync.Mutex{},
    }
}

func (me *EndpointsStruct) StartRock() {
    go func(){
        ticker := time.NewTicker(tickDuration)
        for range ticker.C {
            me.UpdateList()
        }
    }()
}

func (me EndpointsStruct) MemberList() ([]*etcdserverpb.Member, error) {
    resp, err := Request("MemberList")
    if err != nil {
        return nil, err
    }
    r := resp.(*clientv3.MemberListResponse)
    return r.Members, nil
}

func (me EndpointsStruct) Status(endpoint string) (*clientv3.StatusResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
    defer cancel()
    return MaintenanceCli.Status(ctx, endpoint)
}

func (me EndpointsStruct) Health(endpoint string) (latency string, err error) {
    var cli *clientv3.Client
    cli, err = clientv3.New(clientv3.Config{
        Endpoints: []string{endpoint},
        DialTimeout: DialTimeout,
    })
    if err != nil {
        return
    }
    defer cli.Close()
    startTime := time.Now()
    ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
    _, err = cli.Get(ctx, "health")
    cancel()
    if err != nil {
        return
    }
    latency = time.Since(startTime).String()
    return
}

func (me EndpointsStruct) AlermList() (*clientv3.AlarmResponse, error) {
    ctx, cancel := context.WithTimeout(context.Background(), RequestTimeout)
    defer cancel()
    return MaintenanceCli.AlarmList(ctx)
}

func (me *EndpointsStruct) UpdateList() {

    defer func(){
        err := recover()
        if err != nil {
            me.List = make([]EndpointDetail, 0)
        }
    }()

    members, err := me.MemberList()
    if err != nil {
        panic(err)
        mo_log.Logger.Info("get memberlist failed: %s", err.Error())
        members = make([]*etcdserverpb.Member, 0)
    }

    memberHash := make(map[string]EndpointDetail)

    var wg sync.WaitGroup

    endpoints := mo_conf.Conf.Etcd
    for _, endpoint := range endpoints {
        wg.Add(1)
        go func(edpt mo_conf.EtcdConf) {
            defer wg.Done()
            s, e := me.Status(edpt.Addr)
            latency, err := me.Health(edpt.Addr)
            tem := EndpointDetail{}
            if e == nil {
                tem.AddrForDashboard = edpt.Addr
                tem.Version = s.Version
                tem.DbSize = s.DbSize
                tem.Leader = s.Leader
                tem.RaftIndex = s.RaftIndex
                tem.RaftTerm = s.RaftTerm
            }
            tem.IsHealth = err == nil
            if tem.IsHealth {
                tem.Latency = latency
            }
            me.mu.Lock()
            memberHash[edpt.Name] = tem
            me.mu.Unlock()
        }(endpoint)
    }
    wg.Wait()

    for _, mb := range members {
        if detail, ok := memberHash[mb.Name]; ok {
            detail.Member = mb
            me.mu.Lock()
            memberHash[mb.Name] = detail
            me.mu.Unlock()
        } else {
            me.mu.Lock()
            memberHash[mb.Name] = EndpointDetail{
                Member: mb,
            }
            me.mu.Unlock()
        }
    }

    list := make([]EndpointDetail, 0)
    for _, m := range memberHash {
        list = append(list, m)
    }

    me.List = list

}

