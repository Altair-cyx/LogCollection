package etcd

import (
	"context"
	"encoding/json"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

var (
	cli *clientv3.Client
)

// LogEntry 需要收集的日志的配置信息
type LogEntry struct{
	Path string `json:"path"` // 日志存放的路径
	Topic string `json:"topic"` // 日志要发往kafka中的哪个Topic
}



// Init 初始化ETCD的函数
func Init(addr string, timeout time.Duration)(err error){
	cli, err = clientv3.New(clientv3.Config{
		Endpoints:   []string{addr},
		DialTimeout: 5 * timeout,
	})
	if err != nil {
		// handle error!
		fmt.Printf("connect to etcd failed, err:%v\n", err)
		return
	}
	return
}

// GetConf 从etcd中根据key获取配置项
func GetConf(key string)(logEntryConf []*LogEntry,err error){
	// get
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := cli.Get(ctx, key)
	cancel()
	if err != nil {
		fmt.Printf("get from etcd failed, err:%v\n", err)
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("%s:%s\n", ev.Key, ev.Value)
		err = json.Unmarshal(ev.Value,&logEntryConf)
		if err != nil{
			fmt.Printf("unmarshal etcd value failed,err:%v\n",err)
			return
		}
	}
	return
}

func WatchConf(key string, newConfCh chan<- []*LogEntry){
	// 派一个哨兵 一直监视着kongjie的动向这个key的变化（新增、修改、删除）
	ch := cli.Watch(context.Background(),key)
	// 从通道尝试取值（监视的信息）
	for wresp := range ch{
		for _,evt := range wresp.Events{
			fmt.Printf("Type:%v key:%v value:%v\n",evt.Type,string(evt.Kv.Key),string(evt.Kv.Value))
			// 通知taillog.tskMgr
			// 1.先判断操作的类型
			var newConf []*LogEntry
			if evt.Type != clientv3.EventTypeDelete {
				// 如果是删除操作，手动传递一个空的配置项
				err := json.Unmarshal(evt.Kv.Value, &newConf)
				if err != nil{
					fmt.Printf("unmarshal failed,err:%v\n",err)
					continue
				}
			}
			fmt.Printf("get new conf:%v\n",newConf)
			newConfCh <- newConf
		}
	}
}