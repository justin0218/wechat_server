package etcd

import (
	"wechat_server/store"
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

type Wiper struct {
	client *clientv3.Client
}

func NewClient() (*Wiper, error) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(conf.Get().Etcd.Addr, ","),
		DialTimeout: 15 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &Wiper{client: cli}, nil
}

var log = new(store.Log)
var conf = new(store.Config)
var v = fmt.Sprintf("/%s/%s/%s", conf.Get().Etcd.Schema, conf.Get().Etcd.Name, conf.Get().Etcd.Key)

func Register() {
	client, err := NewClient()
	if err != nil {
		log.Get().Error("connect to etcd err:%v", err)
		return
	}
	client.KeepAlive()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	go func() {
		sig := <-ch
		client.UnRegister(conf.Get().Etcd.Name, conf.Get().Etcd.Key)
		if i, ok := sig.(syscall.Signal); ok {
			os.Exit(int(i))
		} else {
			os.Exit(0)
		}
	}()
	log.Get().Info("etcd connect success!!!")
}

func (s *Wiper) GetClient() *clientv3.Client {
	return s.client
}

// Register 注册地址到ETCD组件中 使用 ; 分割
func (s *Wiper) KeepAlive() {
	ticker := time.NewTicker(time.Second * time.Duration(conf.Get().Etcd.Ttl))
	go func() {
		for {
			getResp, err := s.client.Get(context.Background(), v)
			if err != nil {
				log.Get().Error("cli getResp err:%+v\n %v", getResp, err)
				continue
			} else if getResp.Count == 0 {
				leaseResp, err := s.client.Grant(context.Background(), conf.Get().Etcd.Ttl)
				if err != nil {
					log.Get().Error("etcd client.Grant error:%v", err)
					continue
				}
				_, err = s.client.Put(context.Background(), v, conf.Get().Etcd.Key, clientv3.WithLease(leaseResp.ID))
				if err != nil {
					log.Get().Error("etcd s.client.Put error:%v", err)
					continue
				}
				ch, err := s.client.KeepAlive(context.Background(), leaseResp.ID)
				if err != nil {
					log.Get().Error("etcd s.client.KeepAlive err:%v", err)
					continue
				}
				// 清空 keep alive 返回的channel
				go func() {
					for {
						<-ch
					}
				}()
			}
			<-ticker.C
		}
	}()
}

// UnRegister remove service from etcd
func (s *Wiper) UnRegister(name string, addr string) {
	s.client.Delete(context.Background(), v)
}
