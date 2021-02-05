package etcd

import (
	"wechat_server/store"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/resolver"
	"strings"
)

func Discovery(serverName string) *grpc.ClientConn {
	conf := new(store.Config)
	r := NewResolver()
	resolver.Register(r)
	conn, err := grpc.Dial(r.Scheme()+"://"+conf.Get().Etcd.Schema+"/"+serverName, grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)), grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return conn
}

// etcdResolver 解析struct
type etcdResolver struct {
	rawAddr string
	cc      resolver.ClientConn
	client  *clientv3.Client
}

// NewResolver initialize an etcd client
func NewResolver() resolver.Builder {
	etcdWiper, err := NewClient()
	if err != nil {
		log.Get().Error("discovery NewClient err:%v", err)
		return nil
	}
	return &etcdResolver{rawAddr: conf.Get().Etcd.Addr, client: etcdWiper.GetClient()}
}

// Build 构建etcd client
func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r.cc = cc
	go r.watch("/" + target.Scheme + "/" + target.Endpoint + "/")
	return r, nil
}

// Scheme etcd resolve scheme
func (r etcdResolver) Scheme() string {
	return conf.Get().Etcd.Schema
}

// ResolveNow
func (r etcdResolver) ResolveNow(rn resolver.ResolveNowOptions) {

}

// Close closes the resolver
func (r etcdResolver) Close() {

}

// watch 监听resolve列表变化
func (r *etcdResolver) watch(keyPrefix string) {
	var addrList []resolver.Address
	getResp, err := r.client.Get(context.Background(), keyPrefix, clientv3.WithPrefix())
	if err != nil {
		log.Get().Error("discovery etcdResolver watch err:%v", err)
		return
	} else {
		for i := range getResp.Kvs {
			addrList = append(addrList, resolver.Address{Addr: strings.TrimPrefix(string(getResp.Kvs[i].Key), keyPrefix)})
		}
	}
	fmt.Println(addrList, keyPrefix)
	// 新版本etcd去除了NewAddress方法 以UpdateState代替
	r.cc.UpdateState(resolver.State{Addresses: addrList})
	rch := r.client.Watch(context.Background(), keyPrefix, clientv3.WithPrefix())
	for n := range rch {
		for _, ev := range n.Events {
			addr := strings.TrimPrefix(string(ev.Kv.Key), keyPrefix)
			switch ev.Type {
			case mvccpb.PUT:
				if !exist(addrList, addr) {
					addrList = append(addrList, resolver.Address{Addr: addr})
					r.cc.UpdateState(resolver.State{Addresses: addrList})
				}
			case mvccpb.DELETE:
				if s, ok := remove(addrList, addr); ok {
					addrList = s
					r.cc.UpdateState(resolver.State{Addresses: addrList})
				}
			}
			//log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
		}
	}
}

// exist 判断resolve address是否存在
func exist(l []resolver.Address, addr string) bool {
	for i := range l {
		if l[i].Addr == addr {
			return true
		}
	}
	return false
}

// remove 从resolver列表移除
func remove(s []resolver.Address, addr string) ([]resolver.Address, bool) {
	for i := range s {
		if s[i].Addr == addr {
			s[i] = s[len(s)-1]
			return s[:len(s)-1], true
		}
	}
	return nil, false
}
